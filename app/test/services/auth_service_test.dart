import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/services/auth_service.dart';
import 'package:ferna/services/http_client.dart';
import 'package:mockito/mockito.dart';
import '../mocks.mocks.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  group('AuthService.signup', () {
    late MockDio mockDio;
    late AuthService authService;

    setUp(() async {
      // Create a temporary directory for cookie storage
      final tempDir = Directory.systemTemp.createTempSync();
      final cookieDirPath = '${tempDir.path}/cookies_test/';

      // Initialize HttpClient, passing the temp folder so path_provider isn't needed
      await HttpClient.instance.init(
        baseUrl: 'http://example.com',
        testCookieDir: cookieDirPath,
      );

      // Create and assign the mocked Dio
      mockDio = MockDio();
      HttpClient.instance.dio = mockDio;

      // Grab the singleton AuthService (it will use HttpClient.instance.dio)
      authService = AuthService.instance;
    });

    test('returns user_id on success (201 + success:true)', () async {
      // Arrange: stub the mock response with status 201 and {"success": true, "user_id": 42}
      final fakeResponse = Response(
        data: {'success': true, 'user_id': 42},
        statusCode: 201,
        requestOptions: RequestOptions(path: '/auth/local/signup'),
      );
      when(
        mockDio.post(
          '/auth/local/signup',
          data: anyNamed('data'),
        ),
      ).thenAnswer((_) async => fakeResponse);

      // Act:
      final result = await authService.signup(
        email: 'test@example.com',
        password: 'password123',
      );

      // Assert:
      expect(result, 42);

      // The actual implementation sends { "user": ..., "passwd": ... }
      verify(
        mockDio.post(
          '/auth/local/signup',
          data: {'user': 'test@example.com', 'passwd': 'password123'},
        ),
      ).called(1);
    });

    test('throws when server returns success:false', () async {
      // Arrange: server returns 200 but { "success": false, "error": "Invalid data" }
      final fakeResponse = Response(
        data: {'success': false, 'error': 'Invalid data'},
        statusCode: 200,
        requestOptions: RequestOptions(path: '/auth/local/signup'),
      );
      when(
        mockDio.post(
          '/auth/local/signup',
          data: anyNamed('data'),
        ),
      ).thenAnswer((_) async => fakeResponse);

      // Act & Assert: signup() should throw because success is false.
      expect(
        () async => await authService.signup(
          email: 'bad@example.com',
          password: 'short',
        ),
        throwsA(isA<Exception>()),
      );
    });

    test('throws when HTTP status is non-200/201', () async {
      // Arrange: server returns 500 status
      final fakeResponse = Response(
        data: 'Internal Server Error',
        statusCode: 500,
        requestOptions: RequestOptions(path: '/auth/local/signup'),
        statusMessage: 'Internal Server Error',
      );
      when(
        mockDio.post(
          '/auth/local/signup',
          data: anyNamed('data'),
        ),
      ).thenAnswer((_) async => fakeResponse);

      // Act & Assert:
      expect(
        () async => await authService.signup(
          email: 'test@example.com',
          password: 'password123',
        ),
        throwsA(isA<Exception>()),
      );
    });
  });
}
