import 'dart:io';
import 'dart:developer' as developer;

import 'package:dio/dio.dart';
import 'package:cookie_jar/cookie_jar.dart';
import 'package:dio_cookie_manager/dio_cookie_manager.dart';
import 'package:path_provider/path_provider.dart';

class HttpClient {
  HttpClient._();

  static final HttpClient instance = HttpClient._();

  late Dio dio;
  late final PersistCookieJar cookieJar;
  bool _cookieJarInitialized = false;
  String? _currentBaseUrl;

  /// Initialize the HTTP client with base URL and cookie management, testCookieDir is optional directory for storing cookies during tests
  Future<void> init({required String baseUrl, String? testCookieDir}) async {
    developer.log('HttpClient: Initializing with baseUrl: $baseUrl', name: 'ferna.http');
    
    if (!_cookieJarInitialized) {
      // Figure out cookie path
      final cookiesPath = testCookieDir ??
          (await (() async {
            final appDocDir = await getApplicationDocumentsDirectory();
            return '${appDocDir.path}/.cookies/';
          })());

      // Create persistent cookie jar
      cookieJar = PersistCookieJar(
        ignoreExpires: false,
        storage: FileStorage(cookiesPath),
      );
      _cookieJarInitialized = true;
      developer.log('HttpClient: Cookie jar initialized at: $cookiesPath', name: 'ferna.http');
    }

    // Store current base URL
    _currentBaseUrl = baseUrl;

    // Create dio with base url and enhanced configuration
    dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 15),
        receiveTimeout: const Duration(seconds: 15),
        sendTimeout: const Duration(seconds: 15),
        headers: {
          HttpHeaders.contentTypeHeader: 'application/json',
          HttpHeaders.acceptHeader: 'application/json',
        },
        validateStatus: (status) {
          // Accept all status codes to handle them manually at a higher level
          return status != null && status < 500;
        },
      ),
    );

    // Add cookie manager
    dio.interceptors.add(CookieManager(cookieJar));

    // Add request/response interceptor for logging and error handling
    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          developer.log(
            'HttpClient: ${options.method} ${options.uri}',
            name: 'ferna.http',
          );
          handler.next(options);
        },
        onResponse: (response, handler) {
          developer.log(
            'HttpClient: ${response.statusCode} ${response.requestOptions.method} ${response.requestOptions.uri}',
            name: 'ferna.http',
          );
          handler.next(response);
        },
        onError: (error, handler) {
          developer.log(
            'HttpClient: Error ${error.response?.statusCode} ${error.requestOptions.method} ${error.requestOptions.uri} - ${error.message}',
            name: 'ferna.http',
          );
          handler.next(error);
        },
      ),
    );

    developer.log('HttpClient: Initialization complete', name: 'ferna.http');
  }

  /// Clear all cookies (useful for logout)
  Future<void> clearCookies() async {
    developer.log('HttpClient: Clearing all cookies', name: 'ferna.http');
    await cookieJar.deleteAll();
  }

  /// Get current base URL
  String? get baseUrl => _currentBaseUrl;

  /// Check if client is initialized
  bool get isInitialized => _cookieJarInitialized;

  /// Make a POST request
  Future<Response<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) async {
    try {
      return await dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
    } on DioException catch (e) {
      developer.log('HttpClient: POST request failed: ${e.message}', name: 'ferna.http');
      rethrow;
    }
  }

  /// Make a GET request
  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) async {
    try {
      return await dio.get<T>(
        path,
        queryParameters: queryParameters,
        options: options,
      );
    } on DioException catch (e) {
      developer.log('HttpClient: GET request failed: ${e.message}', name: 'ferna.http');
      rethrow;
    }
  }
}
