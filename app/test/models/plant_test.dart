import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/plant.dart';

void main() {
  group('Plant Model', () {
    const testPlantJson = {
      'id': 1,
      'user_id': 123,
      'species_id': 456,
      'nickname': 'My Plant',
      'image_url': 'https://example.com/image.jpg',
      'watering_frequency_days': 7,
      'last_watered_at': '2025-07-01T10:30:00.000Z',
      'note': 'Living room plant',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    const testPlantJsonMinimal = {
      'id': 2,
      'user_id': 123,
      'species_id': 456,
      'watering_frequency_days': 14,
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    test('creates Plant from JSON with all fields', () {
      final plant = Plant.fromJson(testPlantJson);

      expect(plant.id, equals(1));
      expect(plant.userId, equals(123));
      expect(plant.speciesId, equals(456));
      expect(plant.nickname, equals('My Plant'));
      expect(plant.imageUrl, equals('https://example.com/image.jpg'));
      expect(plant.wateringFrequencyDays, equals(7));
      expect(plant.lastWateredAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
      expect(plant.note, equals('Living room plant'));
      expect(plant.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(plant.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('creates Plant from JSON with minimal fields (nullable fields as null)', () {
      final plant = Plant.fromJson(testPlantJsonMinimal);

      expect(plant.id, equals(2));
      expect(plant.userId, equals(123));
      expect(plant.speciesId, equals(456));
      expect(plant.nickname, isNull);
      expect(plant.imageUrl, isNull);
      expect(plant.wateringFrequencyDays, equals(14));
      expect(plant.lastWateredAt, isNull);
      expect(plant.note, isNull);
      expect(plant.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(plant.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('converts Plant to JSON', () {
      final plant = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Plant',
        imageUrl: 'https://example.com/image.jpg',
        wateringFrequencyDays: 7,
        lastWateredAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
        note: 'Living room plant',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final json = plant.toJson();

      expect(json['id'], equals(1));
      expect(json['user_id'], equals(123));
      expect(json['species_id'], equals(456));
      expect(json['nickname'], equals('My Plant'));
      expect(json['image_url'], equals('https://example.com/image.jpg'));
      expect(json['watering_frequency_days'], equals(7));
      expect(json['last_watered_at'], equals('2025-07-01T10:30:00.000Z'));
      expect(json['note'], equals('Living room plant'));
      expect(json['created_at'], equals('2025-06-01T10:30:00.000Z'));
      expect(json['updated_at'], equals('2025-07-01T10:30:00.000Z'));
    });

    test('copyWith creates new instance with updated fields', () {
      final originalPlant = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'Original Name',
        imageUrl: null,
        wateringFrequencyDays: 7,
        lastWateredAt: null,
        note: null,
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final updatedPlant = originalPlant.copyWith(
        nickname: 'Updated Name',
        imageUrl: 'https://example.com/new-image.jpg',
        note: 'New note',
      );

      expect(updatedPlant.id, equals(originalPlant.id));
      expect(updatedPlant.userId, equals(originalPlant.userId));
      expect(updatedPlant.speciesId, equals(originalPlant.speciesId));
      expect(updatedPlant.nickname, equals('Updated Name'));
      expect(updatedPlant.imageUrl, equals('https://example.com/new-image.jpg'));
      expect(updatedPlant.wateringFrequencyDays, equals(originalPlant.wateringFrequencyDays));
      expect(updatedPlant.lastWateredAt, equals(originalPlant.lastWateredAt));
      expect(updatedPlant.note, equals('New note'));
      expect(updatedPlant.createdAt, equals(originalPlant.createdAt));
      expect(updatedPlant.updatedAt, equals(originalPlant.updatedAt));
    });

    test('equality works correctly', () {
      final plant1 = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Plant',
        imageUrl: 'https://example.com/image.jpg',
        wateringFrequencyDays: 7,
        lastWateredAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
        note: 'Living room plant',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final plant2 = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Plant',
        imageUrl: 'https://example.com/image.jpg',
        wateringFrequencyDays: 7,
        lastWateredAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
        note: 'Living room plant',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final plant3 = plant1.copyWith(nickname: 'Different Name');

      expect(plant1, equals(plant2));
      expect(plant1, isNot(equals(plant3)));
      expect(plant1.hashCode, equals(plant2.hashCode));
      expect(plant1.hashCode, isNot(equals(plant3.hashCode)));
    });

    test('toString returns meaningful representation', () {
      final plant = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Plant',
        imageUrl: null,
        wateringFrequencyDays: 7,
        lastWateredAt: null,
        note: null,
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final stringRepresentation = plant.toString();
      
      expect(stringRepresentation, contains('Plant('));
      expect(stringRepresentation, contains('id: 1'));
      expect(stringRepresentation, contains('userId: 123'));
      expect(stringRepresentation, contains('nickname: My Plant'));
    });
  });
}
