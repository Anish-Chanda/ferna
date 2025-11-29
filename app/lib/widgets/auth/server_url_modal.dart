import 'package:flutter/material.dart';
import '../../constants.dart';
import '../global/buttons.dart';
import '../global/inputs.dart';

/// Modal for entering custom server URL
class ServerUrlModal extends StatefulWidget {
  const ServerUrlModal({super.key});

  @override
  State<ServerUrlModal> createState() => _ServerUrlModalState();
}

class _ServerUrlModalState extends State<ServerUrlModal> {
  final _formKey = GlobalKey<FormState>();
  final _serverUrlController = TextEditingController();

  @override
  void dispose() {
    _serverUrlController.dispose();
    super.dispose();
  }

  void _handleSave() {
    if (_formKey.currentState?.validate() ?? false) {
      // TODO: Save to shared preferences later
      Navigator.of(context).pop(_serverUrlController.text);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: AppConstants.surfaceColor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppConstants.radiusLG),
      ),
      child: Padding(
        padding: const EdgeInsets.all(AppConstants.spaceLG),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'Custom Server',
                  style: TextStyle(
                    fontSize: AppConstants.fontSizeLG,
                    fontWeight: AppConstants.fontWeightSemibold,
                    color: AppConstants.textPrimary,
                  ),
                ),
                IconButton(
                  onPressed: () => Navigator.of(context).pop(),
                  icon: const Icon(
                    Icons.close,
                    color: AppConstants.textSecondary,
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: AppConstants.spaceMD),
            
            // Description
            const Text(
              'Enter your self-hosted Ferna server URL to connect to your own instance.',
              style: TextStyle(
                fontSize: AppConstants.fontSizeSM,
                color: AppConstants.textSecondary,
                height: 1.5,
              ),
            ),
            
            const SizedBox(height: AppConstants.spaceLG),
            
            // Form
            Form(
              key: _formKey,
              child: Column(
                children: [
                  ServerUrlTextField(
                    controller: _serverUrlController,
                  ),
                  
                  const SizedBox(height: AppConstants.spaceLG),
                  
                  // Buttons
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () => Navigator.of(context).pop(),
                          style: OutlinedButton.styleFrom(
                            minimumSize: const Size(0, AppConstants.buttonHeightSM),
                          ),
                          child: const Text('Cancel'),
                        ),
                      ),
                      
                      const SizedBox(width: AppConstants.spaceMD),
                      
                      Expanded(
                        child: PrimaryButton(
                          text: 'Save',
                          onPressed: _handleSave,
                          fullWidth: false,
                          height: AppConstants.buttonHeightSM,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Shows the server URL modal and returns the entered URL
Future<String?> showServerUrlModal(BuildContext context) {
  return showDialog<String>(
    context: context,
    builder: (context) => const ServerUrlModal(),
  );
}