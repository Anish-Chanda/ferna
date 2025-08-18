import 'package:ferna/services/http_client.dart';

class AuthService {
  AuthService._();

  static final AuthService instance = AuthService._();

  //Must be called before using any otehr methods of AuthService
  Future<void> initialize({required String serverUrl}) async {
    await HttpClient.instance.init(baseUrl: serverUrl);
  }

  // Sign up a new user. Returns the user_id on success.
  Future<int> signup({required String email, required String password}) async {
    final dio = HttpClient.instance.dio;
    final response = await dio.post(
      '/auth/local/signup',
      data: {'user': email, 'passwd': password},
    );

    if (response.statusCode == 201) {
      final data = response.data as Map<String, dynamic>;
      if (data['success'] == true && data['user_id'] != null) {
        return data['user_id'] as int;
      } else {
        throw Exception('Sign-up failed: ${response.data}');
      }
    } else {
      throw Exception(
        'Sign-up HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  Future<void> login({required String email, required String password}) async {
    final dio = HttpClient.instance.dio;
    final response = await dio.post(
      '/auth/local/login?session=1',
      data: {'user': email, 'passwd': password},
    );

    if (response.statusCode == 200) {
      // CookieManager already saved any Set-Cookie headers to disk.
      return;
    } else {
      throw Exception(
        'Login HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Check if user is currently authenticated by making a test API call
  Future<bool> checkAuthState() async {
    try {
      final dio = HttpClient.instance.dio;
      // Try to access a protected endpoint - the backend should return 401 if not authenticated
      // TODO: add /api/me for this pupose
      final response = await dio.get('/api/plants');

      // If we get a successful response (200-299), user is authenticated
      return response.statusCode != null &&
          response.statusCode! >= 200 &&
          response.statusCode! < 300;
    } catch (e) {
      // Any error (network, 401, 403, etc.) means not authenticated
      return false;
    }
  }

  // Logout user by clearing cookies
  Future<void> logout() async {
    await HttpClient.instance.clearCookies();
  }
}
