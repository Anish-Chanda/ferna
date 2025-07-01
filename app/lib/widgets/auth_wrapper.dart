import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../screens/login_screen.dart';

/// Wrapper that checks authentication state and routes to appropriate screen
class AuthWrapper extends StatelessWidget {
  final Widget authenticatedChild;

  const AuthWrapper({super.key, required this.authenticatedChild});

  @override
  Widget build(BuildContext context) {
    return Consumer<AuthProvider>(
      builder: (context, auth, _) {
        // Show loading indicator while checking auth state
        if (auth.isCheckingAuthState) {
          return const Scaffold(
            body: Center(child: CircularProgressIndicator()),
          );
        }

        // Route based on authentication state
        if (auth.isAuthenticated) {
          log('User is authenticated, navigating to home screen');
          return authenticatedChild;
        } else {
          return const LoginScreen();
        }
      },
    );
  }
}
