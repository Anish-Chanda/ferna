import 'package:flutter/material.dart';

/// Design constants for the Ferna app
class AppConstants {
  AppConstants._();

  // Colors
  static const Color primaryGreen = Color(0xFF12332B);
  static const Color primaryDark = Color(0xFF0A1F1A);
  static const Color primaryLight = Color(0xFF1E5045);
  
  // UI Colors
  static const Color backgroundDark = Color(0xFF1A1A1A);
  static const Color surfaceColor = Color(0xFFFFFFFF);
  static const Color surfaceColorDark = Color(0xFF2C2C2C);
  
  // Text Colors
  static const Color textPrimary = Color(0xFF212121);
  static const Color textSecondary = Color(0xFF757575);
  static const Color textTertiary = Color(0xFF9E9E9E);
  static const Color textOnDark = Color(0xFFFFFFFF);
  static const Color textOnDarkSecondary = Color(0xB3FFFFFF);
  
  // Form Colors
  static const Color inputBorder = Color(0xFFE0E0E0);
  static const Color inputFill = Color(0xFFF5F5F5);
  static const Color inputFocus = Color(0xFF2E7D32);
  
  // Status Colors
  static const Color error = Color(0xFFD32F2F);
  static const Color success = Color(0xFF388E3C);
  static const Color warning = Color(0xFFF57C00);
  static const Color info = Color(0xFF1976D2);
  
  // Social Colors
  static const Color google = Color(0xFF4285F4);
  static const Color facebook = Color(0xFF1877F2);
  
  // Spacing
  static const double spaceXS = 4.0;
  static const double spaceSM = 8.0;
  static const double spaceMD = 16.0;
  static const double spaceLG = 24.0;
  static const double spaceXL = 32.0;
  static const double spaceXXL = 48.0;
  
  // Border Radius
  static const double radiusXS = 4.0;
  static const double radiusSM = 8.0;
  static const double radiusMD = 12.0;
  static const double radiusLG = 16.0;
  static const double radiusXL = 20.0;
  static const double radiusXXL = 24.0;
  
  // Elevation
  static const double elevationSM = 2.0;
  static const double elevationMD = 4.0;
  static const double elevationLG = 8.0;
  static const double elevationXL = 12.0;
  
  // Font Sizes
  static const double fontSizeXS = 12.0;
  static const double fontSizeSM = 14.0;
  static const double fontSizeMD = 16.0;
  static const double fontSizeLG = 18.0;
  static const double fontSizeXL = 20.0;
  static const double fontSizeXXL = 24.0;
  static const double fontSizeDisplay = 32.0;
  static const double fontSizeTitle = 28.0;
  
  // Font Weights
  static const FontWeight fontWeightRegular = FontWeight.w400;
  static const FontWeight fontWeightMedium = FontWeight.w500;
  static const FontWeight fontWeightSemibold = FontWeight.w600;
  static const FontWeight fontWeightBold = FontWeight.w700;
  
  // Input Heights
  static const double inputHeight = 48.0;
  static const double buttonHeight = 48.0;
  static const double buttonHeightSM = 36.0;
  
  // Animation Durations
  static const Duration animationFast = Duration(milliseconds: 200);
  static const Duration animationNormal = Duration(milliseconds: 300);
  static const Duration animationSlow = Duration(milliseconds: 500);
  
  // Layout
  static const double maxContentWidth = 400.0;
  static const double headerHeight = 0.2; // 20% of screen height
  
  // Auth Screen Specific
  static const double backgroundOpacity = 0.6;
  static const double blurRadius = 3.0;
}