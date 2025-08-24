import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/location.dart';

void main() {
  group('Location Model', () {
    const testLocationJson = {
      'id': 1,
      'user_id': 123,
      'name': 'Living Room',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    const testLocationJsonMinimal = {
      'id': 2,
      'user_id': 456,
      'name': 'Kitchen',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    test('creates Location from JSON with all fields', () {
      final location = Location.fromJson(testLocationJson);

      expect(location.id, equals(1));
      expect(location.userId, equals(123));
      expect(location.name, equals('Living Room'));
      expect(location.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(location.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('creates Location from JSON with different values', () {
      final location = Location.fromJson(testLocationJsonMinimal);

      expect(location.id, equals(2));
      expect(location.userId, equals(456));
      expect(location.name, equals('Kitchen'));
      expect(location.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(location.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('converts Location to JSON', () {
      final location = Location(
        id: 1,
        userId: 123,
        name: 'Living Room',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final json = location.toJson();

      expect(json['id'], equals(1));
      expect(json['user_id'], equals(123));
      expect(json['name'], equals('Living Room'));
      expect(json['created_at'], equals('2025-06-01T10:30:00.000Z'));
      expect(json['updated_at'], equals('2025-07-01T10:30:00.000Z'));
    });

    test('copyWith creates new instance with updated fields', () {
      final originalLocation = Location(
        id: 1,
        userId: 123,
        name: 'Original Name',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final updatedLocation = originalLocation.copyWith(
        name: 'Updated Name',
      );

      expect(updatedLocation.id, equals(originalLocation.id));
      expect(updatedLocation.userId, equals(originalLocation.userId));
      expect(updatedLocation.name, equals('Updated Name'));
      expect(updatedLocation.createdAt, equals(originalLocation.createdAt));
      expect(updatedLocation.updatedAt, equals(originalLocation.updatedAt));
    });

    test('equality works correctly', () {
      final location1 = Location(
        id: 1,
        userId: 123,
        name: 'Living Room',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final location2 = Location(
        id: 1,
        userId: 123,
        name: 'Living Room',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final location3 = location1.copyWith(name: 'Different Name');

      expect(location1, equals(location2));
      expect(location1, isNot(equals(location3)));
      expect(location1.hashCode, equals(location2.hashCode));
      expect(location1.hashCode, isNot(equals(location3.hashCode)));
    });

    test('toString returns meaningful representation', () {
      final location = Location(
        id: 1,
        userId: 123,
        name: 'Living Room',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final stringRepresentation = location.toString();
      
      expect(stringRepresentation, contains('Location('));
      expect(stringRepresentation, contains('id: 1'));
      expect(stringRepresentation, contains('userId: 123'));
      expect(stringRepresentation, contains('name: Living Room'));
    });
  });
}
