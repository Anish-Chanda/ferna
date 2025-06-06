import 'package:ferna/services/http_client.dart';
import 'package:flutter/material.dart';
// ignore: unused_import
import 'package:provider/provider.dart';

import '../services/auth_service.dart';
import '../services/storage_service.dart';

class AuthProvider with ChangeNotifier {
  final AuthService _authService = AuthService.instance;
  AuthProvider._();
  bool _isLoading = false;
  bool get isLoading => _isLoading;
  String _serverUrl = '';
  String get serverUrl => _serverUrl;

  // Must call this before using any other AuthProvider methods.
  static Future<AuthProvider> initialize() async {
    final provider = AuthProvider._();

    // Load saved serverUrl or use a default.
    final savedUrl = await StorageService.getServerUrl();
    provider._serverUrl = savedUrl ?? 'http://ferna.local';

    // Initialize AuthService (sets up Dio & loads any persisted cookies).
    await provider._authService.initialize(serverUrl: provider._serverUrl);

    return provider;
  }

  /// Change serverUrl, persist it, and re-initialize Dio so its baseUrl updates.
  Future<void> updateServerUrl(String newUrl) async {
    if (newUrl == _serverUrl) return;

    _serverUrl = newUrl;
    notifyListeners();

    // Persist the new URL:
    await StorageService.saveServerUrl(newUrl);

    // Re-init only Dio’s baseUrl; reuse the same cookieJar
    await HttpClient.instance.init(baseUrl: _serverUrl);

    notifyListeners();
  }

  Future<void> login({required String email, required String password}) async {
    _isLoading = true;
    notifyListeners();

    try {
      await _authService.login(email: email, password: password);
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<int> signUp({required String email, required String password}) async {
    _isLoading = true;
    notifyListeners();

    try {
      final userId = await _authService.signup(
        email: email,
        password: password,
      );
      _isLoading = false;
      notifyListeners();
      return userId;
    } catch (e) {
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }
}
