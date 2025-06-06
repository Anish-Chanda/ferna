import 'package:ferna/providers/auth_provider.dart';
import 'package:ferna/screens/login_screen.dart';
import 'package:ferna/theme/theme_data.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialzie auth provider
  final authProvider = await AuthProvider.initialize();

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider<AuthProvider>.value(value: authProvider),
      ],
      child: FernaApp(),
    ),
  );
}

class FernaApp extends StatelessWidget {
  const FernaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Ferna',
      debugShowCheckedModeBanner: false,
      theme: FernaTheme.light,
      darkTheme: FernaTheme.dark,
      themeMode: ThemeMode.system,
      // TODO: Replace with router
      home: const LoginScreen(),
      routes: {'/home': (_) => const HomeScreen()},
    );
  }
}

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Placeholder(child: Text(":)"));
  }
}
