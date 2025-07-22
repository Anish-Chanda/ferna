import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/plant.dart';
import 'package:ferna/models/plant_display_model.dart';

void main() {
  group('PlantDisplayModel', () {
    late Plant testPlant;
    late DateTime testCreatedAt;
    late DateTime testUpdatedAt;
    late DateTime testLastWateredAt;

    setUp(() {
      testCreatedAt = DateTime.parse('2025-06-01T10:30:00.000Z');
      testUpdatedAt = DateTime.parse('2025-07-01T10:30:00.000Z');
      testLastWateredAt = DateTime.now().subtract(const Duration(days: 3));
      
      testPlant = Plant(
        id: 1,
        userId: 123,
        speciesId: 456,
        nickname: 'My Monstera',
        imageUrl: 'https://example.com/image.jpg',
        wateringFrequencyDays: 7,
        lastWateredAt: testLastWateredAt,
        note: 'Living room',
        createdAt: testCreatedAt,
        updatedAt: testUpdatedAt,
      );
    });

    test('displayName returns nickname when available', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('My Monstera'));
    });

    test('displayName returns species common name when nickname is null', () {
      final plantWithoutNickname = testPlant.copyWith(nickname: null);
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('Monstera Deliciosa'));
    });

    test('displayName returns Unknown Plant when both nickname and species name are null', () {
      final plantWithoutNickname = testPlant.copyWith(nickname: null);
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        speciesCommonName: null,
        location: 'Living room',
      );

      expect(displayModel.displayName, equals('Unknown Plant'));
    });

    test('locationText returns provided location', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Bedroom',
      );

      expect(displayModel.locationText, equals('Bedroom'));
    });

    test('locationText returns No location when location is null', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        speciesCommonName: 'Monstera Deliciosa',
        location: null,
      );

      expect(displayModel.locationText, equals('No location'));
    });

    test('hasCustomNickname returns true when nickname exists', () {
      final displayModel = PlantDisplayModel(
        plant: testPlant,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isTrue);
    });

    test('hasCustomNickname returns false when nickname is null', () {
      final plantWithoutNickname = testPlant.copyWith(nickname: null);
      final displayModel = PlantDisplayModel(
        plant: plantWithoutNickname,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isFalse);
    });

    test('hasCustomNickname returns false when nickname is empty', () {
      final plantWithEmptyNickname = testPlant.copyWith(nickname: '');
      final displayModel = PlantDisplayModel(
        plant: plantWithEmptyNickname,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.hasCustomNickname, isFalse);
    });

    test('daysSinceWatered calculates correctly', () {
      final now = DateTime.now();
      final fiveDaysAgo = now.subtract(const Duration(days: 5));
      final plantWateredFiveDaysAgo = testPlant.copyWith(lastWateredAt: fiveDaysAgo);
      
      final displayModel = PlantDisplayModel(
        plant: plantWateredFiveDaysAgo,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.daysSinceWatered, equals(5));
    });

    test('daysSinceWatered returns 0 when lastWateredAt is null', () {
      final plantNeverWatered = testPlant.copyWith(lastWateredAt: null);
      final displayModel = PlantDisplayModel(
        plant: plantNeverWatered,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.daysSinceWatered, equals(0));
    });

    test('daysOverdue calculates correctly when overdue', () {
      final now = DateTime.now();
      final tenDaysAgo = now.subtract(const Duration(days: 10));
      final plantWateredTenDaysAgo = testPlant.copyWith(
        lastWateredAt: tenDaysAgo,
        wateringFrequencyDays: 7,
      );
      
      final displayModel = PlantDisplayModel(
        plant: plantWateredTenDaysAgo,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.daysOverdue, equals(3)); // 10 - 7 = 3 days overdue
    });

    test('daysOverdue is negative when not overdue', () {
      final now = DateTime.now();
      final threeDaysAgo = now.subtract(const Duration(days: 3));
      final plantWateredThreeDaysAgo = testPlant.copyWith(
        lastWateredAt: threeDaysAgo,
        wateringFrequencyDays: 7,
      );
      
      final displayModel = PlantDisplayModel(
        plant: plantWateredThreeDaysAgo,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.daysOverdue, equals(-4)); // 3 - 7 = -4 days (not overdue)
    });

    test('needsWatering returns true when overdue', () {
      final now = DateTime.now();
      final tenDaysAgo = now.subtract(const Duration(days: 10));
      final plantWateredTenDaysAgo = testPlant.copyWith(
        lastWateredAt: tenDaysAgo,
        wateringFrequencyDays: 7,
      );
      
      final displayModel = PlantDisplayModel(
        plant: plantWateredTenDaysAgo,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.needsWatering, isTrue);
    });

    test('needsWatering returns false when not overdue', () {
      final now = DateTime.now();
      final threeDaysAgo = now.subtract(const Duration(days: 3));
      final plantWateredThreeDaysAgo = testPlant.copyWith(
        lastWateredAt: threeDaysAgo,
        wateringFrequencyDays: 7,
      );
      
      final displayModel = PlantDisplayModel(
        plant: plantWateredThreeDaysAgo,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.needsWatering, isFalse);
    });

    test('nextWateringDate calculates correctly', () {
      final lastWatered = DateTime(2025, 7, 1, 10, 30);
      final plantWithKnownWateringDate = testPlant.copyWith(
        lastWateredAt: lastWatered,
        wateringFrequencyDays: 7,
      );
      
      final displayModel = PlantDisplayModel(
        plant: plantWithKnownWateringDate,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      final expectedNextWatering = DateTime(2025, 7, 8, 10, 30);
      expect(displayModel.nextWateringDate, equals(expectedNextWatering));
    });

    test('nextWateringDate returns null when lastWateredAt is null', () {
      final plantNeverWatered = testPlant.copyWith(lastWateredAt: null);
      final displayModel = PlantDisplayModel(
        plant: plantNeverWatered,
        speciesCommonName: 'Monstera Deliciosa',
        location: 'Living room',
      );

      expect(displayModel.nextWateringDate, isNull);
    });
  });
}
