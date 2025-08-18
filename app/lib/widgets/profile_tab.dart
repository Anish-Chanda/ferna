import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';

class ProfileTab extends StatelessWidget {
  const ProfileTab({super.key});

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header
              Text('Profile', style: theme.textTheme.headlineLarge),
              const SizedBox(height: 24),

              // Profile Card
              _buildProfileCard(context, auth),
              const SizedBox(height: 24),

              // Settings Section
              _buildSectionHeader(context, 'Settings'),
              const SizedBox(height: 16),
              _buildSettingsCard(context, [
                _buildSettingItem(
                  context: context,
                  icon: Icons.notifications_outlined,
                  title: 'Notifications',
                  subtitle: 'Manage your notification preferences',
                  onTap: () {
                    // TODO: Navigate to notifications settings
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Notifications settings coming soon!'),
                      ),
                    );
                  },
                ),
                _buildDivider(context),
                _buildSettingItem(
                  context: context,
                  icon: Icons.palette_outlined,
                  title: 'Theme',
                  subtitle: 'Choose your app appearance',
                  onTap: () {
                    // TODO: Navigate to theme settings
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Theme settings coming soon!'),
                      ),
                    );
                  },
                ),
                _buildDivider(context),
                _buildSettingItem(
                  context: context,
                  icon: Icons.language_outlined,
                  title: 'Language',
                  subtitle: 'Change app language',
                  onTap: () {
                    // TODO: Navigate to language settings
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Language settings coming soon!'),
                      ),
                    );
                  },
                ),
              ]),
              const SizedBox(height: 24),

              // Account Section
              _buildSectionHeader(context, 'Account'),
              const SizedBox(height: 16),
              _buildSettingsCard(context, [
                _buildSettingItem(
                  context: context,
                  icon: Icons.settings,
                  title: 'Server Settings',
                  subtitle: auth.serverUrl,
                  onTap: () {
                    _showServerConfigDialog(context, auth);
                  },
                ),
              ]),
              const SizedBox(height: 24),

              // Logout Button
              _buildLogoutButton(context, auth),
              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProfileCard(BuildContext context, AuthProvider auth) {
    final theme = Theme.of(context);

    return Center(
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              // Profile Avatar
              CircleAvatar(
                radius: 40,
                backgroundColor: theme.chipTheme.backgroundColor,
                child: Icon(
                  Icons.person,
                  size: 40,
                  color: theme.colorScheme.primary,
                ),
              ),
              const SizedBox(height: 16),

              // User Email
              Text(auth.email, style: theme.textTheme.titleLarge),
              const SizedBox(height: 4),

              // User Status
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 4,
                ),
                decoration: BoxDecoration(
                  color: theme.chipTheme.backgroundColor,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  auth.isAuthenticated ? 'Connected' : 'Not connected',
                  style: theme.chipTheme.labelStyle,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSectionHeader(BuildContext context, String title) {
    final theme = Theme.of(context);
    return Text(title, style: theme.textTheme.titleLarge);
  }

  Widget _buildSettingsCard(BuildContext context, List<Widget> children) {
    return Card(child: Column(children: children));
  }

  Widget _buildSettingItem({
    required BuildContext context,
    required IconData icon,
    required String title,
    required String subtitle,
    required VoidCallback onTap,
  }) {
    final theme = Theme.of(context);
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: theme.chipTheme.backgroundColor,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(icon, color: theme.colorScheme.primary, size: 20),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: theme.textTheme.bodyLarge!.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: theme.textTheme.bodyMedium,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            Icon(
              Icons.chevron_right,
              color: theme.colorScheme.onSurface.withOpacity(0.4),
              size: 20,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDivider(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Divider(height: 1, color: theme.dividerTheme.color),
    );
  }

  Widget _buildLogoutButton(BuildContext context, AuthProvider auth) {
    final theme = Theme.of(context);
    final errorColor = theme.colorScheme.error;

    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: () => _showLogoutDialog(context, auth),
        style: ElevatedButton.styleFrom(
          backgroundColor: errorColor.withOpacity(0.1),
          foregroundColor: errorColor,
          elevation: 0,
          padding: const EdgeInsets.symmetric(vertical: 16),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.logout, color: errorColor),
            const SizedBox(width: 8),
            Text(
              'Log Out',
              style: theme.textTheme.bodyLarge!.copyWith(
                fontWeight: FontWeight.w600,
                color: errorColor,
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showLogoutDialog(BuildContext context, AuthProvider auth) {
    final theme = Theme.of(context);

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Log Out'),
          content: const Text('Are you sure you want to log out?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () async {
                Navigator.of(context).pop(); // Close dialog
                await auth.logout();
                // AuthWrapper will automatically handle navigation to login screen
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: theme.colorScheme.error,
                foregroundColor: theme.colorScheme.onError,
              ),
              child: const Text('Log Out'),
            ),
          ],
        );
      },
    );
  }

  void _showServerConfigDialog(BuildContext context, AuthProvider auth) {
    final TextEditingController serverUrlController = TextEditingController(
      text: auth.serverUrl,
    );

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Configure Server'),
          content: TextFormField(
            controller: serverUrlController,
            decoration: InputDecoration(
              labelText: 'Server URL',
              hintText: 'https://ferna.local',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            keyboardType: TextInputType.url,
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () async {
                final newUrl = serverUrlController.text.trim();
                if (newUrl.isNotEmpty && newUrl != auth.serverUrl) {
                  await auth.updateServerUrl(newUrl);
                }
                if (context.mounted) {
                  Navigator.of(context).pop();
                }
              },
              child: const Text('Save'),
            ),
          ],
        );
      },
    );
  }
}
