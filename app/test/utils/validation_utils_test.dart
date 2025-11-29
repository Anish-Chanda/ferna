import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/utils/validation_utils.dart';

void main() {
  group('ValidationUtils.validateServerUrl', () {
    test('should return error for empty or null input', () {
      expect(ValidationUtils.validateServerUrl(null), 'Please enter a server URL');
      expect(ValidationUtils.validateServerUrl(''), 'Please enter a server URL');
    });

    test('should return error for URLs without http/https protocol', () {
      expect(
        ValidationUtils.validateServerUrl('example.com'),
        'URL must start with http:// or https://',
      );
      expect(
        ValidationUtils.validateServerUrl('ftp://example.com'),
        'URL must start with http:// or https://',
      );
    });

    test('should accept valid HTTP URLs', () {
      expect(ValidationUtils.validateServerUrl('http://example.com'), null);
      expect(ValidationUtils.validateServerUrl('https://example.com'), null);
      expect(ValidationUtils.validateServerUrl('https://api.example.com'), null);
      expect(ValidationUtils.validateServerUrl('http://localhost:3000'), null);
      expect(ValidationUtils.validateServerUrl('https://192.168.1.1'), null);
      expect(ValidationUtils.validateServerUrl('http://192.168.1.1:8080'), null);
    });

    test('should reject invalid domain formats', () {
      expect(
        ValidationUtils.validateServerUrl('https://test@gmail.com'),
        'Invalid host format',
      );
      expect(
        ValidationUtils.validateServerUrl('http://user:pass@example.com'),
        'Invalid host format',
      );
    });

    test('should reject malformed URLs', () {
      expect(
        ValidationUtils.validateServerUrl('https://'),
        'URL must contain a valid host/domain',
      );
      expect(
        ValidationUtils.validateServerUrl('http:// '),
        'Invalid URL format',
      );
    });

    test('should validate port numbers correctly', () {
      expect(ValidationUtils.validateServerUrl('http://example.com:80'), null);
      expect(ValidationUtils.validateServerUrl('https://example.com:443'), null);
      expect(ValidationUtils.validateServerUrl('http://example.com:8080'), null);
      expect(
        ValidationUtils.validateServerUrl('http://example.com:99999'),
        'Invalid port number',
      );
    });

    test('should validate IP addresses correctly', () {
      expect(ValidationUtils.validateServerUrl('http://127.0.0.1'), null);
      expect(ValidationUtils.validateServerUrl('https://192.168.1.1'), null);
      expect(ValidationUtils.validateServerUrl('http://[::1]'), null);
      expect(ValidationUtils.validateServerUrl('https://[2001:db8::1]'), null);
    });

    test('should reject invalid IP addresses', () {
      expect(
        ValidationUtils.validateServerUrl('http://256.1.1.1'),
        'Invalid domain format',
      );
      expect(
        ValidationUtils.validateServerUrl('http://192.168.1'),
        'Invalid domain format',
      );
    });
  });

  group('ValidationUtils.validateEmail', () {
    test('should return error for empty or null input', () {
      expect(ValidationUtils.validateEmail(null), 'Please enter your email');
      expect(ValidationUtils.validateEmail(''), 'Please enter your email');
    });

    test('should accept valid email addresses', () {
      expect(ValidationUtils.validateEmail('test@example.com'), null);
      expect(ValidationUtils.validateEmail('user.name@domain.co.uk'), null);
      expect(ValidationUtils.validateEmail('test+tag@example.org'), null);
    });

    test('should reject invalid email addresses', () {
      expect(
        ValidationUtils.validateEmail('invalid-email'),
        'Please enter a valid email address',
      );
      expect(
        ValidationUtils.validateEmail('@example.com'),
        'Please enter a valid email address',
      );
      expect(
        ValidationUtils.validateEmail('test@'),
        'Please enter a valid email address',
      );
    });
  });

  group('ValidationUtils.validatePassword', () {
    test('should return error for empty or null input', () {
      expect(ValidationUtils.validatePassword(null), 'Please enter your password');
      expect(ValidationUtils.validatePassword(''), 'Please enter your password');
    });

    test('should accept any password for login', () {
      expect(ValidationUtils.validatePassword('123'), null);
      expect(ValidationUtils.validatePassword('short'), null);
    });

    test('should enforce minimum length for signup', () {
      expect(
        ValidationUtils.validatePassword('123', isSignUp: true),
        'Password must be at least 6 characters',
      );
      expect(ValidationUtils.validatePassword('123456', isSignUp: true), null);
      expect(ValidationUtils.validatePassword('longerpassword', isSignUp: true), null);
    });
  });
}