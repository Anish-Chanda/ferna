import 'package:ferna/screens/home_screen.dart';
import 'package:flutter/material.dart';
import 'widgets/auth_wrapper.dart';

// Main app navigation widget and handles auth state
class AppNavigator extends StatelessWidget {
  const AppNavigator({super.key});

  @override
  Widget build(BuildContext context) {
    return const AuthWrapper(
      authenticatedChild: HomeScreen(),
    );
  }
}
