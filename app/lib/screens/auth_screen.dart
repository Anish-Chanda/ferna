import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:developer' as developer;
import 'dart:ui';
import 'package:google_fonts/google_fonts.dart';
import '../providers/auth_provider.dart';
import '../constants.dart';
import '../widgets/auth/auth_tab_switcher.dart';
import '../widgets/auth/server_url_modal.dart';
import '../widgets/global/buttons.dart';
import '../widgets/global/inputs.dart';
import '../services/storage_service.dart';
import '../utils/validation_utils.dart';

class AuthScreen extends StatefulWidget {
  const AuthScreen({super.key});

  @override
  State<AuthScreen> createState() => _AuthScreenState();
}

class _AuthScreenState extends State<AuthScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _nameController = TextEditingController();
  final _scrollController = ScrollController();
  
  bool _isSignUp = false;
  bool _isLoading = false;
  String _serverUrl = 'api.fernalabs.com';

  @override
  void initState() {
    super.initState();
    developer.log('AuthScreen: Initializing authentication screen', name: 'ferna.screens');
    _loadCurrentServerUrl();
  }

  /// Load the current server URL to display in the UI
  Future<void> _loadCurrentServerUrl() async {
    try {
      if (!mounted) return;
      
      // Get the actual server URL from storage service
      final baseUrl = await StorageService.getServerUrl();
      final displayUrl = ValidationUtils.getDisplayUrl(baseUrl);
      
      if (mounted) {
        setState(() {
          _serverUrl = displayUrl;
        });
      }
    } catch (e) {
      developer.log('AuthScreen: Failed to load server URL: $e', name: 'ferna.screens');
      // Fallback to default display
      if (mounted) {
        setState(() {
          _serverUrl = 'api.fernalabs.com';
        });
      }
    }
  }

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _nameController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _handleAuth() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);
    
    final authProvider = Provider.of<AuthProvider>(context, listen: false);
    final email = _emailController.text.trim();
    final password = _passwordController.text;
    final name = _nameController.text.trim();

    developer.log(
      'AuthScreen: Attempting ${_isSignUp ? 'sign up' : 'sign in'} for email: $email', 
      name: 'ferna.screens'
    );

    bool success;
    if (_isSignUp) {
      success = await authProvider.signUp(
        email: email, 
        password: password,
        name: name, // Pass name directly, validation will handle empty case
      );
    } else {
      success = await authProvider.signIn(email: email, password: password);
    }

    setState(() => _isLoading = false);

    if (!success && mounted) {
      final errorMessage = authProvider.lastError ?? '${_isSignUp ? 'Sign up' : 'Sign in'} failed. Please try again.';
      developer.log('AuthScreen: Authentication failed: $errorMessage', name: 'ferna.screens');
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(errorMessage),
          backgroundColor: AppConstants.error,
        ),
      );
    } else {
      developer.log('AuthScreen: Authentication successful', name: 'ferna.screens');
    }
  }

  Future<void> _showServerModal() async {
    final result = await showServerUrlModal(context);
    if (result != null && result.isNotEmpty) {
      if (!mounted) return;
      
      try {
        final authProvider = Provider.of<AuthProvider>(context, listen: false);
        await authProvider.updateServerUrl(result);
        
        if (mounted) {
          setState(() {
            _serverUrl = ValidationUtils.getHostFromUrl(result);
          });
        }
        
        developer.log('AuthScreen: Server URL updated to: $result', name: 'ferna.screens');
      } catch (e) {
        developer.log('AuthScreen: Failed to update server URL: $e', name: 'ferna.screens');
        
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Failed to update server URL: $e'),
              backgroundColor: AppConstants.error,
            ),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    developer.log('AuthScreen: Building authentication screen widget', name: 'ferna.screens');
    
    final screenSize = MediaQuery.of(context).size;
    // More flexible header height - minimum 160px, but not more than 22% of screen
    final headerHeight = (screenSize.height * 0.22).clamp(160.0, screenSize.height * 0.25);

    return Scaffold(
      body: Stack(
        children: [
          // Background Image with Blur
          Container(
            height: double.infinity,
            width: double.infinity,
            decoration: const BoxDecoration(
              image: DecorationImage(
                image: AssetImage('assets/auth_background_leaves.jpg'),
                fit: BoxFit.cover,
              ),
            ),
            child: BackdropFilter(
              filter: ImageFilter.blur(
                sigmaX: AppConstants.blurRadius,
                sigmaY: AppConstants.blurRadius,
              ),
              child: Container(
                color: AppConstants.backgroundDark.withValues(alpha: AppConstants.backgroundOpacity),
              ),
            ),
          ),

          // Content
          SafeArea(
            child: Column(
              children: [
                // Header Section (20% of screen)
                Container(
                  height: headerHeight,
                  width: double.infinity,
                  padding: const EdgeInsets.only(
                    left: AppConstants.spaceLG,
                    right: AppConstants.spaceLG,
                    top: AppConstants.spaceMD,
                    bottom: AppConstants.spaceMD,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // App Logo/Icon and Title
                          Row(
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              Image.asset(
                                'assets/icon/ferna_sprout_only.png',
                                height: 40,
                                width: 40,
                              ),
                              const SizedBox(width: AppConstants.spaceSM),
                              Padding(
                                padding: const EdgeInsets.only(top: 10),
                                child: Text(
                                  'Ferna',
                                  style: GoogleFonts.nunito(
                                    fontSize: 32,
                                    fontWeight: FontWeight.w600, // Semibold
                                    color: AppConstants.textOnDark,
                                    height: 1,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: AppConstants.spaceLG),
                          
                          // Main Title
                          const Text(
                            'Get Started now',
                            style: TextStyle(
                              fontSize: 28,
                              fontWeight: AppConstants.fontWeightBold,
                              color: AppConstants.textOnDark,
                              height: 1.1,
                            ),
                          ),
                          
                          const SizedBox(height: AppConstants.spaceSM),
                          
                          // Subtitle
                          const Text(
                            'Create an account or log in',
                            style: TextStyle(
                              fontSize: AppConstants.fontSizeMD,
                              fontWeight: AppConstants.fontWeightRegular,
                              color: AppConstants.textOnDarkSecondary,
                              height: 1.3,
                            ),
                          ),
                        ],
                      ),
                      
                      // Photo Attribution
                      const Text(
                        'Photo by Thimo van Leeuwen on Unsplash',
                        style: TextStyle(
                          fontSize: AppConstants.fontSizeXS,
                          color: AppConstants.textOnDarkSecondary,
                        ),
                      ),
                    ],
                  ),
                ),

                // Bottom Section - Auth Form with Scrolling
                Expanded(
                  child: Container(
                    width: double.infinity,
                    decoration: const BoxDecoration(
                      color: AppConstants.surfaceColor,
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(AppConstants.radiusXXL),
                        topRight: Radius.circular(AppConstants.radiusXXL),
                      ),
                    ),
                    child: SingleChildScrollView(
                      controller: _scrollController,
                      physics: const ClampingScrollPhysics(),
                      child: Padding(
                        padding: const EdgeInsets.all(AppConstants.spaceLG),
                        child: ConstrainedBox(
                          constraints: const BoxConstraints(maxWidth: AppConstants.maxContentWidth),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              // Tab Switcher
                              AuthTabSwitcher(
                                isSignUp: _isSignUp,
                                onTabChanged: (isSignUp) {
                                  setState(() {
                                    _isSignUp = isSignUp;
                                    // Clear name field when switching to login
                                    if (!isSignUp) {
                                      _nameController.clear();
                                    }
                                  });
                                  developer.log(
                                    'AuthScreen: Switched to ${isSignUp ? 'sign up' : 'sign in'} mode', 
                                    name: 'ferna.screens'
                                  );
                                },
                              ),

                              const SizedBox(height: AppConstants.spaceLG),

                              // Auth Form
                              Form(
                                key: _formKey,
                                child: Column(
                                  children: [
                                    // Name Field (only for signup)
                                    if (_isSignUp) ...[
                                      NameTextField(
                                        controller: _nameController,
                                      ),
                                      const SizedBox(height: AppConstants.spaceLG),
                                    ],

                                    // Email Field
                                    EmailTextField(
                                      controller: _emailController,
                                    ),

                                    const SizedBox(height: AppConstants.spaceLG),

                                    // Password Field
                                    PasswordTextField(
                                      controller: _passwordController,
                                      isSignUp: _isSignUp,
                                    ),

                                    // Forgot Password (only for login)
                                    if (!_isSignUp) ...[
                                      const SizedBox(height: AppConstants.spaceSM),
                                      Align(
                                        alignment: Alignment.centerRight,
                                        child: TextButton(
                                          onPressed: () {
                                            developer.log('AuthScreen: Forgot password tapped', name: 'ferna.screens');
                                          },
                                          child: const Text(
                                            'Forgot Password ?',
                                            style: TextStyle(
                                              fontSize: AppConstants.fontSizeSM,
                                              color: AppConstants.textSecondary,
                                            ),
                                          ),
                                        ),
                                      ),
                                      const SizedBox(height: AppConstants.spaceSM),
                                    ],

                                    if (_isSignUp)
                                      const SizedBox(height: AppConstants.spaceMD),

                                    // Auth Button
                                    PrimaryButton(
                                      text: _isSignUp ? 'Sign Up' : 'Log In',
                                      onPressed: _handleAuth,
                                      isLoading: _isLoading,
                                      backgroundColor: _isSignUp 
                                          ? AppConstants.primaryGreen 
                                          : null,
                                    ),

                                    const SizedBox(height: AppConstants.spaceLG),

                                    // Divider
                                    const Row(
                                      children: [
                                        Expanded(child: Divider()),
                                        Padding(
                                          padding: EdgeInsets.symmetric(horizontal: AppConstants.spaceMD),
                                          child: Text(
                                            'Or',
                                            style: TextStyle(
                                              fontSize: AppConstants.fontSizeSM,
                                              color: AppConstants.textTertiary,
                                            ),
                                          ),
                                        ),
                                        Expanded(child: Divider()),
                                      ],
                                    ),

                                    const SizedBox(height: AppConstants.spaceLG),

                                    // Social Login Buttons
                                    GoogleLoginButton(
                                      onPressed: () {
                                        developer.log('AuthScreen: Google login tapped', name: 'ferna.screens');
                                      },
                                    ),

                                    const SizedBox(height: AppConstants.spaceSM),

                                    FacebookLoginButton(
                                      onPressed: () {
                                        developer.log('AuthScreen: Facebook login tapped', name: 'ferna.screens');
                                      },
                                    ),

                                    const SizedBox(height: AppConstants.spaceLG),

                                    // Server Selection
                                    GestureDetector(
                                      onTap: _showServerModal,
                                      child: Container(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: AppConstants.spaceMD,
                                          vertical: AppConstants.spaceSM,
                                        ),
                                        decoration: BoxDecoration(
                                          border: Border.all(color: AppConstants.inputBorder),
                                          borderRadius: BorderRadius.circular(AppConstants.radiusSM),
                                        ),
                                        child: Row(
                                          mainAxisAlignment: MainAxisAlignment.center,
                                          children: [
                                            Text(
                                              'Logging in on: $_serverUrl',
                                              style: const TextStyle(
                                                fontSize: AppConstants.fontSizeSM,
                                                color: AppConstants.primaryGreen,
                                                fontWeight: AppConstants.fontWeightMedium,
                                              ),
                                            ),
                                            const SizedBox(width: AppConstants.spaceSM),
                                            const Icon(
                                              Icons.keyboard_arrow_down,
                                              size: 16,
                                              color: AppConstants.textTertiary,
                                            ),
                                          ],
                                        ),
                                      ),
                                    ),

                                    // Extra spacing for keyboard
                                    SizedBox(height: MediaQuery.of(context).viewInsets.bottom + AppConstants.spaceLG),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
