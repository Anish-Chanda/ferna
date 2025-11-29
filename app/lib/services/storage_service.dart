import 'package:shared_preferences/shared_preferences.dart';

const _kServerUrlKey = 'ferna_server_url';

class StorageService {
  static const String _productionUrl = 'https://api.fernalabs.com';

  static Future<SharedPreferences> _prefs() async =>
      await SharedPreferences.getInstance();

  /// Get the default server URL
  static String getDefaultServerUrl() => _productionUrl;

  /// Save the server URL for future app launches
  static Future<void> saveServerUrl(String url) async {
    final prefs = await _prefs();
    await prefs.setString(_kServerUrlKey, url);
  }

  /// Read the persisted server URL, or return default if none was saved
  static Future<String> getServerUrl() async {
    final prefs = await _prefs();
    final savedUrl = prefs.getString(_kServerUrlKey);
    return savedUrl ?? getDefaultServerUrl();
  }

  /// Clear the saved server URL (will revert to default)
  static Future<void> clearServerUrl() async {
    final prefs = await _prefs();
    await prefs.remove(_kServerUrlKey);
  }
}
