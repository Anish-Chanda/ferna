import 'package:flutter/foundation.dart';
import 'dart:developer' as developer;
import '../services/auth_service.dart';

enum AuthStatus {
  unknown,
  authenticated,
  unauthenticated,
}

/// Authentication provider that manages user authentication state
/// and automatically notifies listeners of state changes
class AuthProvider with ChangeNotifier {
  AuthStatus _status = AuthStatus.unknown;
  String? _userId;
  String? _userEmail;
  String? _userName;
  bool _isLoading = false;
  String? _lastError;

  // Getters
  AuthStatus get status => _status;
  String? get userId => _userId;
  String? get userEmail => _userEmail;
  String? get userName => _userName;
  bool get isLoading => _isLoading;
  String? get lastError => _lastError;
  bool get isAuthenticated => _status == AuthStatus.authenticated;

  /// Initialize the auth provider and check current auth state
  Future<void> initialize() async {
    developer.log('AuthProvider: Initializing authentication state', name: 'ferna.auth');
    
    try {
      // Initialize auth service
      await AuthService.instance.initialize();
      
      // Check if user is already authenticated
      await _checkAuthState();
    } catch (e) {
      developer.log('AuthProvider: Initialization error: $e', name: 'ferna.auth');
      _setError('Failed to initialize authentication');
    }
  }

  /// Check current authentication state by fetching user info
  Future<void> _checkAuthState() async {
    developer.log('AuthProvider: Checking authentication state', name: 'ferna.auth');
    
    try {
      final userInfo = await AuthService.instance.getUserInfo();
      
      if (userInfo != null) {
        _status = AuthStatus.authenticated;
        _userEmail = userInfo.email;
        _userName = userInfo.name;
        _userId = userInfo.email; // Use email as ID for now
        _clearError();
        developer.log('AuthProvider: User is authenticated: ${userInfo.email}', name: 'ferna.auth');
      } else {
        _status = AuthStatus.unauthenticated;
        _clearUserData();
        developer.log('AuthProvider: User is not authenticated', name: 'ferna.auth');
      }
    } catch (e) {
      developer.log('AuthProvider: Error checking auth state: $e', name: 'ferna.auth');
      _status = AuthStatus.unauthenticated;
      _clearUserData();
    }
    
    notifyListeners();
  }

  /// Sign in user with email and password
  Future<bool> signIn({required String email, required String password}) async {
    developer.log('AuthProvider: Attempting to sign in user: $email', name: 'ferna.auth');
    
    _setLoading(true);
    _clearError();
    
    try {
      final result = await AuthService.instance.signIn(
        email: email,
        password: password,
      );
      
      if (result.success) {
        // Refresh auth state to get user info
        await _checkAuthState();
        developer.log('AuthProvider: Sign in successful', name: 'ferna.auth');
        return true;
      } else {
        _setError(result.message ?? 'Sign in failed');
        developer.log('AuthProvider: Sign in failed: ${result.message}', name: 'ferna.auth');
        return false;
      }
    } catch (e) {
      _setError('Sign in failed: $e');
      developer.log('AuthProvider: Sign in error: $e', name: 'ferna.auth');
      return false;
    } finally {
      _setLoading(false);
    }
  }

  /// Sign up user with email, password and name
  Future<bool> signUp({
    required String email, 
    required String password,
    required String name, // Changed from optional to required
  }) async {
    developer.log('AuthProvider: Attempting to sign up user: $email', name: 'ferna.auth');
    
    _setLoading(true);
    _clearError();
    
    try {
      final result = await AuthService.instance.signUp(
        email: email,
        password: password,
        name: name,
      );
      
      if (result.success) {
        // After successful signup, sign in the user
        developer.log('AuthProvider: Signup successful, signing in user', name: 'ferna.auth');
        return await signIn(email: email, password: password);
      } else {
        _setError(result.message ?? 'Sign up failed');
        developer.log('AuthProvider: Sign up failed: ${result.message}', name: 'ferna.auth');
        return false;
      }
    } catch (e) {
      _setError('Sign up failed: $e');
      developer.log('AuthProvider: Sign up error: $e', name: 'ferna.auth');
      return false;
    } finally {
      _setLoading(false);
    }
  }

  /// Sign out user
  Future<void> signOut() async {
    developer.log('AuthProvider: Signing out user', name: 'ferna.auth');
    
    _setLoading(true);
    
    try {
      await AuthService.instance.signOut();
      
      _status = AuthStatus.unauthenticated;
      _clearUserData();
      _clearError();
      
      developer.log('AuthProvider: User signed out successfully', name: 'ferna.auth');
    } catch (e) {
      developer.log('AuthProvider: Sign out error: $e', name: 'ferna.auth');
      // Force sign out even if there's an error
      _status = AuthStatus.unauthenticated;
      _clearUserData();
    } finally {
      _setLoading(false);
      notifyListeners();
    }
  }

  /// Refresh authentication state
  Future<void> refreshAuth() async {
    developer.log('AuthProvider: Refreshing authentication state', name: 'ferna.auth');
    await _checkAuthState();
  }

  /// Update server URL and refresh auth state
  Future<void> updateServerUrl(String newUrl) async {
    developer.log('AuthProvider: Updating server URL to: $newUrl', name: 'ferna.auth');
    
    try {
      await AuthService.instance.updateServerUrl(newUrl);
      // After changing server, user might not be authenticated anymore
      await _checkAuthState();
    } catch (e) {
      developer.log('AuthProvider: Failed to update server URL: $e', name: 'ferna.auth');
      _setError('Failed to update server URL');
    }
  }

  /// Helper method to set loading state
  void _setLoading(bool loading) {
    if (_isLoading != loading) {
      _isLoading = loading;
      notifyListeners();
    }
  }

  /// Helper method to set error state
  void _setError(String error) {
    _lastError = error;
    notifyListeners();
  }

  /// Helper method to clear error state
  void _clearError() {
    if (_lastError != null) {
      _lastError = null;
      notifyListeners();
    }
  }

  /// Helper method to clear user data
  void _clearUserData() {
    _userId = null;
    _userEmail = null;
    _userName = null;
  }
}