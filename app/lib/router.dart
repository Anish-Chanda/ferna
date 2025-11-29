//Main app navigation widget that handles auth state

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:developer' as developer;
import 'providers/auth_provider.dart';
import 'screens/auth_screen.dart';
import 'screens/home_screen.dart';

/// Main navigation widget that manages app routing based on authentication state
class AppNavigator extends StatelessWidget {
  const AppNavigator({super.key});

  @override
  Widget build(BuildContext context) {
    developer.log('AppNavigator: Building app navigator', name: 'ferna.router');
    
    return Consumer<AuthProvider>(
      builder: (context, authProvider, child) {
        developer.log(
          'AppNavigator: Auth status changed to ${authProvider.status}', 
          name: 'ferna.router'
        );
        
        // Show loading screen while auth state is being determined
        if (authProvider.status == AuthStatus.unknown) {
          developer.log('AppNavigator: Showing loading screen', name: 'ferna.router');
          return const _LoadingScreen();
        }
        
        // Navigate based on authentication status
        return AuthWrapper(authProvider: authProvider);
      },
    );
  }
}

/// Wrapper widget that determines which screen to show based on auth state
class AuthWrapper extends StatelessWidget {
  const AuthWrapper({
    super.key,
    required this.authProvider,
  });

  final AuthProvider authProvider;

  @override
  Widget build(BuildContext context) {
    developer.log(
      'AuthWrapper: Rendering screen for auth status: ${authProvider.status}', 
      name: 'ferna.router'
    );
    
    switch (authProvider.status) {
      case AuthStatus.authenticated:
        developer.log('AuthWrapper: User is authenticated, showing home screen', name: 'ferna.router');
        return const HomeScreen();
      
      case AuthStatus.unauthenticated:
        developer.log('AuthWrapper: User is not authenticated, showing auth screen', name: 'ferna.router');
        return const AuthScreen();
      
      case AuthStatus.unknown:
        developer.log('AuthWrapper: Auth status unknown, showing loading screen', name: 'ferna.router');
        return const _LoadingScreen();
    }
  }
}

/// Loading screen shown while authentication state is being determined
class _LoadingScreen extends StatelessWidget {
  const _LoadingScreen();

  @override
  Widget build(BuildContext context) {
    developer.log('_LoadingScreen: Building loading screen', name: 'ferna.router');
    
    return const Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.eco,
              size: 60,
              color: Colors.green,
            ),
            SizedBox(height: 16),
            Text(
              'Ferna',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: Colors.green,
              ),
            ),
            SizedBox(height: 32),
            CircularProgressIndicator(
              color: Colors.green,
            ),
            SizedBox(height: 16),
            Text(
              'Loading...',
              style: TextStyle(
                color: Colors.grey,
              ),
            ),
          ],
        ),
      ),
    );
  }
}