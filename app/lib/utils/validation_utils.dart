/// Validation utilities for form inputs
class ValidationUtils {
  ValidationUtils._();

  /// Validates a URL according to RFC 3986 format for server endpoints
  /// Rejects URLs with userinfo (like https://test@gmail.com)
  static String? validateServerUrl(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please enter a server URL';
    }

    // Must start with http:// or https://
    if (!value.startsWith('http://') && !value.startsWith('https://')) {
      return 'URL must start with http:// or https://';
    }

    // Check for userinfo in URL (anything with @ before the domain)
    // This catches cases like https://test@gmail.com
    if (value.contains('@')) {
      return 'Invalid host format';
    }

    // Check for spaces in URL (malformed URLs)
    if (value.contains(' ')) {
      return 'Invalid URL format';
    }

    // Try to parse the URI
    final uri = Uri.tryParse(value);
    if (uri == null) {
      return 'Invalid URL format';
    }

    // Check required components
    if (uri.scheme.isEmpty || !['http', 'https'].contains(uri.scheme)) {
      return 'URL must use http or https protocol';
    }

    if (uri.host.isEmpty) {
      return 'URL must contain a valid host/domain';
    }

    // Additional check for userinfo in parsed URI
    if (uri.userInfo.isNotEmpty) {
      return 'Invalid host format';
    }

    // Check for valid domain format
    if (!_isValidDomain(uri.host)) {
      return 'Invalid domain format';
    }

    // Check port is valid if provided
    if (uri.hasPort && (uri.port < 1 || uri.port > 65535)) {
      return 'Invalid port number';
    }

    return null; // Valid URL
  }

  /// Helper function to validate domain format
  static bool _isValidDomain(String host) {
    // If it looks like an IP address (contains only digits and dots), validate as IP
    if (_looksLikeIPv4(host)) {
      return _isValidIPv4(host);
    }

    // Check if it's IPv6
    if (_isValidIPv6(host)) {
      return true;
    }

    // Check if it's a valid domain name
    return _isValidDomainName(host);
  }

  /// Checks if a string looks like an IPv4 address (digits and dots only)
  static bool _looksLikeIPv4(String host) {
    return RegExp(r'^[0-9.]+$').hasMatch(host);
  }

  /// Validates IPv4 address
  static bool _isValidIPv4(String ip) {
    final parts = ip.split('.');
    if (parts.length != 4) return false;

    for (final part in parts) {
      // Check if part is a valid number
      final num = int.tryParse(part);
      if (num == null) return false;
      
      // Check range 0-255
      if (num < 0 || num > 255) return false;
      
      // Reject leading zeros (except for "0")
      if (part.length > 1 && part.startsWith('0')) {
        return false;
      }
    }
    return true;
  }

  /// Validates IPv6 address (basic check)
  static bool _isValidIPv6(String ip) {
    try {
      // Remove brackets if present
      final cleanIp = ip.startsWith('[') && ip.endsWith(']') 
          ? ip.substring(1, ip.length - 1) 
          : ip;
      
      Uri.parseIPv6Address(cleanIp);
      return true;
    } catch (e) {
      return false;
    }
  }

  /// Validates domain name format
  static bool _isValidDomainName(String domain) {
    // Basic checks
    if (domain.isEmpty || domain.length > 253) return false;
    if (domain.startsWith('.') || domain.endsWith('.')) return false;
    if (domain.startsWith('-') || domain.endsWith('-')) return false;
    
    // Split into labels
    final labels = domain.split('.');
    
    // Need at least one label for localhost, or multiple for proper domains
    if (labels.isEmpty) return false;
    
    // Allow localhost as special case
    if (domain == 'localhost') return true;
    
    // For other domains, require at least one dot
    if (labels.length < 2) return false;
    
    // Check each label
    for (final label in labels) {
      if (label.isEmpty || label.length > 63) return false;
      if (label.startsWith('-') || label.endsWith('-')) return false;
      
      // Check valid characters (letters, numbers, hyphens)
      if (!RegExp(r'^[a-zA-Z0-9-]+$').hasMatch(label)) {
        return false;
      }
    }
    
    return true;
  }

  /// Validates email format
  static String? validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please enter your email';
    }
    
    final emailRegex = RegExp(
      r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
    );
    
    if (!emailRegex.hasMatch(value)) {
      return 'Please enter a valid email address';
    }
    
    return null;
  }

  /// Validates password format
  static String? validatePassword(String? value, {bool isSignUp = false}) {
    if (value == null || value.isEmpty) {
      return 'Please enter your password';
    }
    
    if (isSignUp && value.length < 6) {
      return 'Password must be at least 6 characters';
    }
    
    return null;
  }

  /// Validates name format
  static String? validateName(String? value, {bool isRequired = false}) {
    if (value == null || value.isEmpty) {
      return isRequired ? 'Please enter your name' : null;
    }
    
    // Trim whitespace
    value = value.trim();
    
    if (value.length < 2) {
      return 'Name must be at least 2 characters';
    }
    
    if (value.length > 100) {
      return 'Name must be less than 100 characters';
    }
    
    // Allow letters, numbers, spaces, common punctuation
    final nameRegex = RegExp(r"^[a-zA-Z0-9\s\-'.]+$");
    if (!nameRegex.hasMatch(value)) {
      return 'Name contains invalid characters';
    }
    
    return null;
  }

  /// Parse a URL and return a user-friendly display string
  /// Returns the host and port if non-standard
  static String getDisplayUrl(String url) {
    try {
      final uri = Uri.parse(url);
      return uri.host + (uri.hasPort && uri.port != 80 && uri.port != 443 ? ':${uri.port}' : '');
    } catch (e) {
      // If parsing fails, try to extract a meaningful part
      if (url.startsWith('http://') || url.startsWith('https://')) {
        final withoutProtocol = url.replaceFirst(RegExp(r'^https?://'), '');
        return withoutProtocol.split('/').first;
      }
      return url;
    }
  }

  /// Extract hostname from URL for display purposes
  /// This is a safe method that won't throw exceptions
  static String getHostFromUrl(String url) {
    try {
      final uri = Uri.parse(url);
      return uri.host;
    } catch (e) {
      return url; // Return original if parsing fails
    }
  }
}