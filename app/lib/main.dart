import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:developer' as developer;
import 'providers/auth_provider.dart';
import 'router.dart';
import 'theme.dart';

void main() {
  developer.log('main: Starting Ferna app', name: 'ferna.main');
  runApp(const FernaApp());
}

class FernaApp extends StatelessWidget {
  const FernaApp({super.key});

  @override
  Widget build(BuildContext context) {
    developer.log('FernaApp: Building main app widget', name: 'ferna.main');
    
    return MultiProvider(
      providers: [
        ChangeNotifierProvider<AuthProvider>(
          create: (context) {
            developer.log('FernaApp: Creating AuthProvider instance', name: 'ferna.main');
            final authProvider = AuthProvider();
            // Initialize the auth provider asynchronously
            authProvider.initialize();
            return authProvider;
          },
        ),
      ],
      child: MaterialApp(
        title: 'Ferna',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.lightTheme,
        home: const AppNavigator(),
      ),
    );
  }
}