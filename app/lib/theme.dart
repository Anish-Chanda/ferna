import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'constants.dart';

/// App theme configuration for Ferna
class AppTheme {
  AppTheme._();

  /// Light theme data
  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      fontFamily: GoogleFonts.inter().fontFamily,
      colorScheme: _lightColorScheme,
      textTheme: _textTheme,
      appBarTheme: _appBarTheme,
      elevatedButtonTheme: _elevatedButtonTheme,
      outlinedButtonTheme: _outlinedButtonTheme,
      inputDecorationTheme: _inputDecorationTheme,
      cardTheme: _cardTheme,
      checkboxTheme: _checkboxTheme,
      dividerTheme: _dividerTheme,
      scaffoldBackgroundColor: AppConstants.surfaceColor,
    );
  }

  /// Light color scheme
  static const ColorScheme _lightColorScheme = ColorScheme(
    brightness: Brightness.light,
    primary: AppConstants.primaryGreen,
    onPrimary: AppConstants.textOnDark,
    secondary: AppConstants.primaryLight,
    onSecondary: AppConstants.textOnDark,
    error: AppConstants.error,
    onError: AppConstants.textOnDark,
    surface: AppConstants.surfaceColor,
    onSurface: AppConstants.textPrimary,
    outline: AppConstants.inputBorder,
  );

  /// Text theme configuration
  static const TextTheme _textTheme = TextTheme(
    // Display styles
    displayLarge: TextStyle(
      fontSize: AppConstants.fontSizeDisplay,
      fontWeight: AppConstants.fontWeightBold,
      color: AppConstants.textPrimary,
      height: 1.2,
    ),
    displayMedium: TextStyle(
      fontSize: AppConstants.fontSizeTitle,
      fontWeight: AppConstants.fontWeightBold,
      color: AppConstants.textPrimary,
      height: 1.3,
    ),
    
    // Headline styles
    headlineLarge: TextStyle(
      fontSize: AppConstants.fontSizeXXL,
      fontWeight: AppConstants.fontWeightBold,
      color: AppConstants.textPrimary,
      height: 1.3,
    ),
    headlineMedium: TextStyle(
      fontSize: AppConstants.fontSizeXL,
      fontWeight: AppConstants.fontWeightSemibold,
      color: AppConstants.textPrimary,
      height: 1.4,
    ),
    
    // Title styles
    titleLarge: TextStyle(
      fontSize: AppConstants.fontSizeLG,
      fontWeight: AppConstants.fontWeightSemibold,
      color: AppConstants.textPrimary,
      height: 1.4,
    ),
    titleMedium: TextStyle(
      fontSize: AppConstants.fontSizeMD,
      fontWeight: AppConstants.fontWeightMedium,
      color: AppConstants.textPrimary,
      height: 1.4,
    ),
    
    // Body styles
    bodyLarge: TextStyle(
      fontSize: AppConstants.fontSizeMD,
      fontWeight: AppConstants.fontWeightRegular,
      color: AppConstants.textPrimary,
      height: 1.5,
    ),
    bodyMedium: TextStyle(
      fontSize: AppConstants.fontSizeSM,
      fontWeight: AppConstants.fontWeightRegular,
      color: AppConstants.textSecondary,
      height: 1.5,
    ),
    bodySmall: TextStyle(
      fontSize: AppConstants.fontSizeXS,
      fontWeight: AppConstants.fontWeightRegular,
      color: AppConstants.textTertiary,
      height: 1.5,
    ),
    
    // Label styles
    labelLarge: TextStyle(
      fontSize: AppConstants.fontSizeSM,
      fontWeight: AppConstants.fontWeightMedium,
      color: AppConstants.textPrimary,
      height: 1.4,
    ),
    labelMedium: TextStyle(
      fontSize: AppConstants.fontSizeXS,
      fontWeight: AppConstants.fontWeightMedium,
      color: AppConstants.textSecondary,
      height: 1.4,
    ),
  );

  /// App bar theme
  static const AppBarTheme _appBarTheme = AppBarTheme(
    backgroundColor: AppConstants.surfaceColor,
    foregroundColor: AppConstants.textPrimary,
    elevation: 0,
    centerTitle: true,
    titleTextStyle: TextStyle(
      fontSize: AppConstants.fontSizeLG,
      fontWeight: AppConstants.fontWeightSemibold,
      color: AppConstants.textPrimary,
    ),
  );

  /// Elevated button theme
  static final ElevatedButtonThemeData _elevatedButtonTheme = ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      backgroundColor: AppConstants.primaryGreen,
      foregroundColor: AppConstants.textOnDark,
      minimumSize: const Size(double.infinity, AppConstants.buttonHeight),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      ),
      elevation: AppConstants.elevationSM,
      textStyle: const TextStyle(
        fontSize: AppConstants.fontSizeMD,
        fontWeight: AppConstants.fontWeightMedium,
      ),
    ),
  );

  /// Outlined button theme
  static final OutlinedButtonThemeData _outlinedButtonTheme = OutlinedButtonThemeData(
    style: OutlinedButton.styleFrom(
      foregroundColor: AppConstants.primaryGreen,
      minimumSize: const Size(double.infinity, AppConstants.buttonHeight),
      side: const BorderSide(color: AppConstants.inputBorder),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      ),
      textStyle: const TextStyle(
        fontSize: AppConstants.fontSizeMD,
        fontWeight: AppConstants.fontWeightMedium,
      ),
    ),
  );

  /// Input decoration theme
  static final InputDecorationTheme _inputDecorationTheme = InputDecorationTheme(
    filled: true,
    fillColor: AppConstants.inputFill,
    border: OutlineInputBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      borderSide: const BorderSide(color: AppConstants.inputBorder),
    ),
    enabledBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      borderSide: const BorderSide(color: AppConstants.inputBorder),
    ),
    focusedBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      borderSide: const BorderSide(color: AppConstants.inputFocus, width: 2),
    ),
    errorBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      borderSide: const BorderSide(color: AppConstants.error),
    ),
    contentPadding: const EdgeInsets.symmetric(
      horizontal: AppConstants.spaceMD,
      vertical: AppConstants.spaceMD,
    ),
    labelStyle: const TextStyle(
      color: AppConstants.textSecondary,
      fontSize: AppConstants.fontSizeMD,
    ),
    hintStyle: const TextStyle(
      color: AppConstants.textTertiary,
      fontSize: AppConstants.fontSizeMD,
    ),
  );

  /// Card theme
  static final CardThemeData _cardTheme = CardThemeData(
    color: AppConstants.surfaceColor,
    elevation: AppConstants.elevationSM,
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusLG),
    ),
  );

  /// Checkbox theme
  static final CheckboxThemeData _checkboxTheme = CheckboxThemeData(
    fillColor: WidgetStateProperty.resolveWith((states) {
      if (states.contains(WidgetState.selected)) {
        return AppConstants.primaryGreen;
      }
      return AppConstants.inputFill;
    }),
    checkColor: WidgetStateProperty.all(AppConstants.textOnDark),
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(AppConstants.radiusXS),
    ),
  );

  /// Divider theme
  static const DividerThemeData _dividerTheme = DividerThemeData(
    color: AppConstants.inputBorder,
    thickness: 1,
    space: AppConstants.spaceLG,
  );
}