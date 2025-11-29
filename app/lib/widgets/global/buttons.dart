import 'package:flutter/material.dart';
import '../../constants.dart';

/// Primary button widget for main actions
class PrimaryButton extends StatelessWidget {
  const PrimaryButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.isLoading = false,
    this.enabled = true,
    this.fullWidth = true,
    this.height = AppConstants.buttonHeight,
    this.backgroundColor,
  });

  final String text;
  final VoidCallback? onPressed;
  final bool isLoading;
  final bool enabled;
  final bool fullWidth;
  final double height;
  final Color? backgroundColor;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: fullWidth ? double.infinity : null,
      height: height,
      child: ElevatedButton(
        onPressed: enabled && !isLoading ? onPressed : null,
        style: backgroundColor != null
            ? ElevatedButton.styleFrom(backgroundColor: backgroundColor)
            : null,
        child: isLoading
            ? const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(
                    AppConstants.textOnDark,
                  ),
                ),
              )
            : Text(text),
      ),
    );
  }
}

/// Social login button with icon and text
class SocialButton extends StatelessWidget {
  const SocialButton({
    super.key,
    required this.text,
    required this.onPressed,
    required this.icon,
    this.backgroundColor = AppConstants.surfaceColor,
    this.textColor = AppConstants.textPrimary,
    this.borderColor = AppConstants.inputBorder,
    this.isLoading = false,
  });

  final String text;
  final VoidCallback? onPressed;
  final Widget icon;
  final Color backgroundColor;
  final Color textColor;
  final Color borderColor;
  final bool isLoading;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: AppConstants.buttonHeight,
      child: OutlinedButton(
        onPressed: isLoading ? null : onPressed,
        style: OutlinedButton.styleFrom(
          backgroundColor: backgroundColor,
          foregroundColor: textColor,
          side: BorderSide(color: borderColor),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppConstants.radiusMD),
          ),
        ),
        child: isLoading
            ? const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(
                    AppConstants.textSecondary,
                  ),
                ),
              )
            : Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  icon,
                  const SizedBox(width: AppConstants.spaceMD),
                  Text(
                    text,
                    style: TextStyle(
                      fontSize: AppConstants.fontSizeMD,
                      fontWeight: AppConstants.fontWeightMedium,
                      color: textColor,
                    ),
                  ),
                ],
              ),
      ),
    );
  }
}

/// Google login button
class GoogleLoginButton extends StatelessWidget {
  const GoogleLoginButton({
    super.key,
    required this.onPressed,
    this.isLoading = false,
  });

  final VoidCallback? onPressed;
  final bool isLoading;

  @override
  Widget build(BuildContext context) {
    return SocialButton(
      text: 'Continue with Google',
      onPressed: onPressed,
      isLoading: isLoading,
      icon: Container(
        width: 20,
        height: 20,
        decoration: const BoxDecoration(
          image: DecorationImage(
            image: NetworkImage(
              'https://developers.google.com/identity/images/g-logo.png',
            ),
            fit: BoxFit.contain,
          ),
        ),
      ),
    );
  }
}

/// Facebook login button
class FacebookLoginButton extends StatelessWidget {
  const FacebookLoginButton({
    super.key,
    required this.onPressed,
    this.isLoading = false,
  });

  final VoidCallback? onPressed;
  final bool isLoading;

  @override
  Widget build(BuildContext context) {
    return SocialButton(
      text: 'Continue with Facebook',
      onPressed: onPressed,
      isLoading: isLoading,
      icon: const Icon(
        Icons.facebook,
        size: 20,
        color: AppConstants.facebook,
      ),
    );
  }
}