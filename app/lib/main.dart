import 'package:ferna/providers/auth_provider.dart';
import 'package:ferna/router.dart';
import 'package:ferna/theme/theme_data.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/date_symbol_data_local.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize date formatting for table_calendar
  await initializeDateFormatting();

  // Initialize auth provider
  final authProvider = await AuthProvider.initialize();

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider<AuthProvider>.value(value: authProvider),
      ],
      child: const FernaApp(),
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
      home: const AppNavigator(),
    );
  }
}
