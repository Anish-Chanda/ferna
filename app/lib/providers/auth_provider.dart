import 'package:flutter/material.dart';
// ignore: unused_import
import 'package:provider/provider.dart';

class AuthProvider with ChangeNotifier {
  bool _isLoading = false;
  bool get isLoading => _isLoading;

  Future<void> login({
    required String email,
    required String password,
    required String serverUrl,
  }) async {
    _isLoading = true;
    notifyListeners();

    // TODO: Replace with http call
    await Future.delayed(const Duration(seconds: 2));

    // On success:
    _isLoading = false;
    notifyListeners();
  }

  Future<void> signUp({
    required String email,
    required String password,
    required String serverUrl,
  }) async {
    _isLoading = true;
    notifyListeners();

    // TODO: Replace with HTTP call.
    await Future.delayed(const Duration(seconds: 2));

    _isLoading = false;
    notifyListeners();
  }
}