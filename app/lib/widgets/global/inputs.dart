import 'package:flutter/material.dart';
import '../../constants.dart';
import '../../utils/validation_utils.dart';

/// Custom text input field so we can have consistent styling
class CustomTextField extends StatefulWidget {
  const CustomTextField({
    super.key,
    required this.label,
    this.hint,
    this.controller,
    this.obscureText = false,
    this.keyboardType = TextInputType.text,
    this.validator,
    this.onChanged,
    this.enabled = true,
    this.prefixIcon,
    this.suffixIcon,
    this.maxLines = 1,
  });

  final String label;
  final String? hint;
  final TextEditingController? controller;
  final bool obscureText;
  final TextInputType keyboardType;
  final String? Function(String?)? validator;
  final void Function(String)? onChanged;
  final bool enabled;
  final Widget? prefixIcon;
  final Widget? suffixIcon;
  final int maxLines;

  @override
  State<CustomTextField> createState() => _CustomTextFieldState();
}

class _CustomTextFieldState extends State<CustomTextField> {
  late bool _obscureText;
  late FocusNode _focusNode;

  @override
  void initState() {
    super.initState();
    _obscureText = widget.obscureText;
    _focusNode = FocusNode();
  }

  @override
  void dispose() {
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          widget.label,
          style: const TextStyle(
            fontSize: AppConstants.fontSizeSM,
            fontWeight: AppConstants.fontWeightMedium,
            color: AppConstants.textSecondary,
          ),
        ),
        const SizedBox(height: AppConstants.spaceSM),
        TextFormField(
          controller: widget.controller,
          focusNode: _focusNode,
          obscureText: _obscureText,
          keyboardType: widget.keyboardType,
          validator: widget.validator,
          onChanged: widget.onChanged,
          enabled: widget.enabled,
          maxLines: widget.maxLines,
          style: const TextStyle(
            fontSize: AppConstants.fontSizeMD,
            color: AppConstants.textPrimary,
          ),
          decoration: InputDecoration(
            hintText: widget.hint ?? widget.label,
            prefixIcon: widget.prefixIcon,
            suffixIcon: widget.obscureText
                ? IconButton(
                    icon: Icon(
                      _obscureText ? Icons.visibility_off : Icons.visibility,
                      color: AppConstants.textTertiary,
                    ),
                    onPressed: () {
                      setState(() {
                        _obscureText = !_obscureText;
                      });
                    },
                  )
                : widget.suffixIcon,
          ),
        ),
      ],
    );
  }
}

/// Email input field with validation
class EmailTextField extends StatelessWidget {
  const EmailTextField({
    super.key,
    this.controller,
    this.onChanged,
  });

  final TextEditingController? controller;
  final void Function(String)? onChanged;

  @override
  Widget build(BuildContext context) {
    return CustomTextField(
      label: 'Email',
      hint: 'Enter your email',
      controller: controller,
      keyboardType: TextInputType.emailAddress,
      onChanged: onChanged,
      validator: ValidationUtils.validateEmail,
    );
  }
}

/// Password input field with validation
class PasswordTextField extends StatelessWidget {
  const PasswordTextField({
    super.key,
    this.controller,
    this.onChanged,
    this.isSignUp = false,
  });

  final TextEditingController? controller;
  final void Function(String)? onChanged;
  final bool isSignUp;

  @override
  Widget build(BuildContext context) {
    return CustomTextField(
      label: 'Password',
      hint: 'Enter your password',
      controller: controller,
      obscureText: true,
      onChanged: onChanged,
      validator: (value) => ValidationUtils.validatePassword(value, isSignUp: isSignUp),
    );
  }
}

/// Name input field for signup
class NameTextField extends StatelessWidget {
  const NameTextField({
    super.key,
    this.controller,
    this.onChanged,
    this.isRequired = true, // Default to required for signup
  });

  final TextEditingController? controller;
  final void Function(String)? onChanged;
  final bool isRequired;

  @override
  Widget build(BuildContext context) {
    return CustomTextField(
      label: 'Full Name',
      hint: 'Enter your full name',
      controller: controller,
      keyboardType: TextInputType.name,
      onChanged: onChanged,
      validator: (value) => ValidationUtils.validateName(value, isRequired: isRequired),
      prefixIcon: const Icon(
        Icons.person_outline,
        color: AppConstants.textSecondary,
        size: 20,
      ),
    );
  }
}

/// Server URL input field
class ServerUrlTextField extends StatelessWidget {
  const ServerUrlTextField({
    super.key,
    this.controller,
    this.onChanged,
  });

  final TextEditingController? controller;
  final void Function(String)? onChanged;

  @override
  Widget build(BuildContext context) {
    return CustomTextField(
      label: 'Server URL',
      hint: 'https://your-server.com',
      controller: controller,
      keyboardType: TextInputType.url,
      onChanged: onChanged,
      validator: ValidationUtils.validateServerUrl,
    );
  }
}