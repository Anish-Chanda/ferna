import 'dart:convert';
import 'dart:developer' as developer;
import 'package:dio/dio.dart';
import '../services/http_client.dart';
import '../services/storage_service.dart';

/// Authentication result model
class AuthResult {
  final bool success;
  final String? message;
  final String? userId;
  final Map<String, dynamic>? userData;

  AuthResult({
    required this.success,
    this.message,
    this.userId,
    this.userData,
  });
}

/// User information model
class UserInfo {
  final String email;
  final String? name;
  final String? audience;

  UserInfo({
    required this.email,
    this.name,
    this.audience,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) {
    return UserInfo(
      email: json['name'] ?? '', // go-pkgz/auth returns email in 'name' field for local provider, TODO: we might need to adjust this based on other providers before v1
      name: json['display_name'],
      audience: json['aud'],
    );
  }
}

/// This service is responsible for making auth related HTTP requests using http_client
/// and managing authentication state
class AuthService {
  static final AuthService instance = AuthService._();
  AuthService._();

  static const String _audience = 'ferna-mobile';

  /// Initialize the auth service with proper base URL
  Future<void> initialize() async {
    try {
      final serverUrl = await StorageService.getServerUrl();
      await HttpClient.instance.init(baseUrl: serverUrl);
      developer.log('AuthService: Initialized with server URL: $serverUrl', name: 'ferna.auth');
    } catch (e) {
      developer.log('AuthService: Initialization failed: $e', name: 'ferna.auth');
      rethrow;
    }
  }

  /// Sign up a new user with email, password and name
  Future<AuthResult> signUp({
    required String email,
    required String password,
    required String name, // Changed from optional to required
  }) async {
    try {
      developer.log('AuthService: Attempting signup for email: $email', name: 'ferna.auth');

      // Get timezone info
      final timezone = DateTime.now().timeZoneName;
      
      final requestData = {
        'email': email,
        'password': password,
        'full_name': name.trim(), // Use the provided name directly
        'timezone': timezone,
      };
      
      developer.log('AuthService: Signup request data: $requestData', name: 'ferna.auth');

      final response = await HttpClient.instance.post(
        '/api/auth/signup',
        data: requestData,
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final responseData = response.data;
        developer.log('AuthService: Signup successful, response: $responseData', name: 'ferna.auth');
        
        // Handle both Map and other response types
        String? message;
        String? userId;
        
        if (responseData is Map<String, dynamic>) {
          message = responseData['message']?.toString();
          userId = responseData['user_id']?.toString();
        }
        
        return AuthResult(
          success: true,
          message: message ?? 'Signup successful',
          userId: userId,
          userData: responseData is Map<String, dynamic> ? responseData : null,
        );
      } else {
        String errorMessage = 'Signup failed';
        if (response.data != null) {
          if (response.data is Map<String, dynamic>) {
            errorMessage = response.data['message']?.toString() ?? 
                          response.data['error']?.toString() ?? 
                          'Signup failed';
          } else if (response.data is String) {
            errorMessage = response.data;
          }
        }
        developer.log('AuthService: Signup failed: $errorMessage', name: 'ferna.auth');
        
        return AuthResult(
          success: false,
          message: errorMessage,
        );
      }
    } on DioException catch (e) {
      final errorMessage = _handleDioError(e, 'signup');
      return AuthResult(
        success: false,
        message: errorMessage,
      );
    } catch (e) {
      developer.log('AuthService: Unexpected signup error: $e', name: 'ferna.auth');
      developer.log('AuthService: Error type: ${e.runtimeType}', name: 'ferna.auth');
      developer.log('AuthService: Stack trace: ${StackTrace.current}', name: 'ferna.auth');
      return AuthResult(
        success: false,
        message: 'An unexpected error occurred: ${e.toString()}',
      );
    }
  }

  /// Sign in an existing user
  Future<AuthResult> signIn({
    required String email,
    required String password,
  }) async {
    try {
      developer.log('AuthService: Attempting signin for email: $email', name: 'ferna.auth');

      final response = await HttpClient.instance.post(
        '/auth/local/login',
        data: {
          'user': email,
          'passwd': password,
          'aud': _audience, // Required audience parameter
        },
      );

      if (response.statusCode == 200) {
        developer.log('AuthService: Signin successful', name: 'ferna.auth');
        
        return AuthResult(
          success: true,
          message: 'Login successful',
          userData: response.data,
        );
      } else {
        final errorMessage = response.data?['error'] ?? 'Login failed';
        developer.log('AuthService: Signin failed: $errorMessage', name: 'ferna.auth');
        
        return AuthResult(
          success: false,
          message: errorMessage,
        );
      }
    } on DioException catch (e) {
      final errorMessage = _handleDioError(e, 'signin');
      return AuthResult(
        success: false,
        message: errorMessage,
      );
    } catch (e) {
      developer.log('AuthService: Unexpected signin error: $e', name: 'ferna.auth');
      return AuthResult(
        success: false,
        message: 'An unexpected error occurred',
      );
    }
  }

  /// Sign out the current user
  Future<AuthResult> signOut() async {
    try {
      developer.log('AuthService: Attempting signout', name: 'ferna.auth');

      final response = await HttpClient.instance.post('/auth/logout');

      // Clear cookies regardless of response status
      await HttpClient.instance.clearCookies();

      if (response.statusCode == 200) {
        developer.log('AuthService: Signout successful', name: 'ferna.auth');
        
        return AuthResult(
          success: true,
          message: 'Logout successful',
        );
      } else {
        // Even if server logout fails, we cleared cookies locally
        developer.log('AuthService: Server logout failed but local cookies cleared', name: 'ferna.auth');
        
        return AuthResult(
          success: true,
          message: 'Logged out locally',
        );
      }
    } catch (e) {
      // Even if logout fails, clear cookies locally
      await HttpClient.instance.clearCookies();
      developer.log('AuthService: Logout error, but cookies cleared: $e', name: 'ferna.auth');
      
      return AuthResult(
        success: true,
        message: 'Logged out locally',
      );
    }
  }

  /// Get current user information
  Future<UserInfo?> getUserInfo() async {
    try {
      developer.log('AuthService: Fetching user info', name: 'ferna.auth');

      final response = await HttpClient.instance.get('/auth/user');

      if (response.statusCode == 200) {
        final userData = response.data as Map<String, dynamic>;
        developer.log('AuthService: User info retrieved successfully', name: 'ferna.auth');
        
        return UserInfo.fromJson(userData);
      } else {
        developer.log('AuthService: Failed to get user info: ${response.statusCode}', name: 'ferna.auth');
        return null;
      }
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        developer.log('AuthService: User not authenticated', name: 'ferna.auth');
        return null;
      }
      developer.log('AuthService: Error getting user info: ${e.message}', name: 'ferna.auth');
      return null;
    } catch (e) {
      developer.log('AuthService: Unexpected error getting user info: $e', name: 'ferna.auth');
      return null;
    }
  }

  /// Check if user is currently authenticated
  Future<bool> isAuthenticated() async {
    final userInfo = await getUserInfo();
    return userInfo != null;
  }

  /// Update server URL and reinitialize HTTP client
  Future<void> updateServerUrl(String newUrl) async {
    try {
      await StorageService.saveServerUrl(newUrl);
      await HttpClient.instance.init(baseUrl: newUrl);
      developer.log('AuthService: Server URL updated to: $newUrl', name: 'ferna.auth');
    } catch (e) {
      developer.log('AuthService: Failed to update server URL: $e', name: 'ferna.auth');
      rethrow;
    }
  }

  /// Handle Dio errors and return user-friendly messages
  String _handleDioError(DioException e, String operation) {
    developer.log('AuthService: DioException during $operation: ${e.message}', name: 'ferna.auth');
    developer.log('AuthService: Response status: ${e.response?.statusCode}', name: 'ferna.auth');
    developer.log('AuthService: Response data: ${e.response?.data}', name: 'ferna.auth');
    
    if (e.response?.statusCode != null) {
      final statusCode = e.response!.statusCode!;
      final responseData = e.response?.data;
      
      // Try to extract error message from response
      String? serverMessage;
      if (responseData is Map<String, dynamic>) {
        serverMessage = responseData['error']?.toString() ?? responseData['message']?.toString();
      } else if (responseData is String) {
        // Handle cases where response is a plain string
        try {
          final decoded = jsonDecode(responseData);
          if (decoded is Map<String, dynamic>) {
            serverMessage = decoded['error']?.toString() ?? decoded['message']?.toString();
          }
        } catch (_) {
          serverMessage = responseData;
        }
      }
      
      switch (statusCode) {
        case 400:
          return serverMessage ?? 'Invalid request. Please check your input.';
        case 401:
          return serverMessage ?? 'Invalid credentials. Please try again.';
        case 403:
          return serverMessage ?? 'Access denied.';
        case 404:
          return 'Service not found. Please check your server URL.';
        case 409:
          return serverMessage ?? 'Email already exists. Please use a different email.';
        case 422:
          return serverMessage ?? 'Invalid input data. Please check your information.';
        case 429:
          return 'Too many requests. Please try again later.';
        case 500:
        case 502:
        case 503:
        case 504:
          return 'Server error. Please try again later.';
        default:
          return serverMessage ?? 'Request failed with status $statusCode';
      }
    }
    
    // Network errors
    if (e.type == DioExceptionType.connectionTimeout ||
        e.type == DioExceptionType.receiveTimeout ||
        e.type == DioExceptionType.sendTimeout) {
      return 'Connection timeout. Please check your internet connection.';
    }
    
    if (e.type == DioExceptionType.connectionError) {
      return 'Connection failed. Please check your internet connection and server URL.';
    }
    
    return 'Network error. Please try again.';
  }
}
