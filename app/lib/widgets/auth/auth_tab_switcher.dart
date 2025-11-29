import 'package:flutter/material.dart';
import '../../constants.dart';

/// Tab switcher widget for Login/Sign up toggle
class AuthTabSwitcher extends StatelessWidget {
  const AuthTabSwitcher({
    super.key,
    required this.isSignUp,
    required this.onTabChanged,
  });

  final bool isSignUp;
  final void Function(bool isSignUp) onTabChanged;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppConstants.inputFill,
        borderRadius: BorderRadius.circular(AppConstants.radiusMD),
      ),
      padding: const EdgeInsets.all(AppConstants.spaceXS),
      child: Row(
        children: [
          Expanded(
            child: _TabButton(
              text: 'Log In',
              isSelected: !isSignUp,
              onTap: () => onTabChanged(false),
            ),
          ),
          Expanded(
            child: _TabButton(
              text: 'Sign up',
              isSelected: isSignUp,
              onTap: () => onTabChanged(true),
            ),
          ),
        ],
      ),
    );
  }
}

/// Individual tab button
class _TabButton extends StatelessWidget {
  const _TabButton({
    required this.text,
    required this.isSelected,
    required this.onTap,
  });

  final String text;
  final bool isSelected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: AppConstants.animationNormal,
        curve: Curves.easeInOut,
        height: 40, // Fixed height to prevent flicker
        padding: const EdgeInsets.symmetric(
          vertical: AppConstants.spaceSM,
          horizontal: AppConstants.spaceMD,
        ),
        decoration: BoxDecoration(
          color: isSelected ? AppConstants.surfaceColor : Colors.transparent,
          borderRadius: BorderRadius.circular(AppConstants.radiusSM),
          boxShadow: isSelected
              ? [
                  BoxShadow(
                    color: AppConstants.textTertiary.withValues(alpha: 0.1),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ]
              : null,
        ),
        child: Center(
          child: Text(
            text,
            style: TextStyle(
              fontSize: AppConstants.fontSizeMD,
              fontWeight: isSelected
                  ? AppConstants.fontWeightMedium
                  : AppConstants.fontWeightRegular,
              color: isSelected
                  ? AppConstants.textPrimary
                  : AppConstants.textSecondary,
            ),
          ),
        ),
      ),
    );
  }
}