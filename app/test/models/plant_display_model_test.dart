import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/plant.dart';
import 'package:ferna/models/species.dart';
import 'package:ferna/models/plant_display_model.dart';

void main() {
  group('PlantDisplayModel', () {
    late Plant testPlant;
    late Species testSpecies;
    late DateTime testCreatedAt;
    late DateTime testUpdatedAt;

    setUp(() {
      testCreatedAt = DateTime.parse('2025-06-01T10:30:00.000Z');
      testUpdatedAt = DateTime.parse('2025-07-01T10:30:00.000Z');
      
      testSpecies = Species(
        id: 456,
        commonName: 'Monstera Deliciosa',
        scientificName: 'Monstera deliciosa',
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      
      testPlant = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Monstera',
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: null,
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
    });

    test('displayName returns nickname when available', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('My Monstera'));
    });

    test('displayName returns species common name when nickname is null', () {
      final plantWithoutNickname = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: null, // Explicitly null
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: null,
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('Monstera Deliciosa'));
    });

    test('displayName returns Unknown Plant when both nickname and species are null', () {
      final plantWithoutNickname = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: null, // Explicitly null
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: null,
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        species: null,
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('Unknown Plant'));
    });

    test('locationText returns provided location', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Bedroom',
      );

      expect(displayModel.locationText, equals('Bedroom'));
    });

    test('locationText returns No location when location is null', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: null,
      );

      expect(displayModel.locationText, equals('No location'));
    });

    test('hasCustomNickname returns true when nickname exists', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isTrue);
    });

    test('hasCustomNickname returns false when nickname is null', () {
      final plantWithoutNickname = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: null, // Explicitly null
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: null,
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isFalse);
    });

    test('hasCustomNickname returns false when nickname is empty', () {
      final plantWithEmptyNickname = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: '', // Explicitly empty
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: null,
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithEmptyNickname,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isFalse);
    });

    test('daysSinceWatered returns 0 (placeholder implementation)', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      // This is a placeholder until care event tracking is implemented
      expect(displayModel.daysSinceWatered, equals(0));
    });

    test('wateringFrequencyDays returns plant override when available', () {
      final plantWithOverride = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Monstera',
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: 10, // Override value
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithOverride,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.wateringFrequencyDays, equals(10));
    });

    test('wateringFrequencyDays returns species default when no plant override', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      expect(displayModel.wateringFrequencyDays, equals(7));
    });

    test('wateringFrequencyDays returns default 7 when no species or override', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: null,
        location: 'Living room',
      );

      expect(displayModel.wateringFrequencyDays, equals(7));
    });

    test('daysOverdue calculates based on frequency and days since watered', () {
      final plantWithOverride = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Monstera',
        imageUrl: 'https://example.com/image.jpg',
        notes: 'Living room plant',
        waterIntervalDaysOverride: 7, // Override value
        fertilizerIntervalDaysOverride: null,
        locationId: 1,
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
      final displayModel = PlantDisplayModel(
        plant: plantWithOverride,
        species: testSpecies,
        location: 'Living room',
      );

      // Since daysSinceWatered is 0 and wateringFrequencyDays is 7
      expect(displayModel.daysOverdue, equals(-7)); // 0 - 7 = -7 (not overdue)
    });

    test('needsWatering returns false with current placeholder implementation', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      // Since daysSinceWatered is always 0, needsWatering will be false
      expect(displayModel.needsWatering, isFalse);
    });

    test('nextWateringDate returns null (placeholder implementation)', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        species: testSpecies,
        location: 'Living room',
      );

      // This is a placeholder until care event tracking is implemented
      expect(displayModel.nextWateringDate, isNull);
    });
  });
}
