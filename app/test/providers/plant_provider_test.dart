import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/providers/plant_provider.dart';
import 'package:ferna/models/plant.dart';
import 'package:ferna/services/plant_service.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';

import 'plant_provider_test.mocks.dart';

@GenerateMocks([PlantService])
void main() {
  group('PlantProvider', () {
    late MockPlantService mockPlantService;
    late PlantProvider plantProvider;

    final testPlant1 = Plant(
      id: 1,
      userId: 123,
      speciesId: 456,
      nickname: 'My Plant',
      imageUrl: 'https://example.com/image.jpg',
      notes: 'Living room plant',
      waterIntervalDaysOverride: 7,
      fertilizerIntervalDaysOverride: null,
      locationId: 1,
      createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
      updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
    );

    final testPlant2 = Plant(
      id: 2,
      userId: 123,
      speciesId: 789,
      nickname: 'Second Plant',
      imageUrl: null,
      notes: null,
      waterIntervalDaysOverride: 10,
      fertilizerIntervalDaysOverride: null,
      locationId: null,
      createdAt: DateTime.parse('2025-06-02T10:30:00.000Z'),
      updatedAt: DateTime.parse('2025-07-02T10:30:00.000Z'),
    );

    setUp(() {
      mockPlantService = MockPlantService();
      plantProvider = PlantProvider(plantService: mockPlantService);
    });

    group('fetchPlants', () {
      test('successfully fetches plants and updates state', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1, testPlant2]);

        // Act
        await plantProvider.fetchPlants();

        // Assert
        expect(plantProvider.isLoading, isFalse);
        expect(plantProvider.error, isNull);
        expect(plantProvider.plants, hasLength(2));
        expect(plantProvider.plants.first.id, equals(1));
        expect(plantProvider.plants.last.id, equals(2));

        verify(mockPlantService.getUserPlants(limit: 20, offset: 0)).called(1);
      });

      test('handles error during fetch and updates state', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenThrow(Exception('Network error'));

        // Act
        await plantProvider.fetchPlants();

        // Assert
        expect(plantProvider.isLoading, isFalse);
        expect(plantProvider.error, contains('Network error'));
        expect(plantProvider.plants, isEmpty);

        verify(mockPlantService.getUserPlants(limit: 20, offset: 0)).called(1);
      });

      test('sets loading state correctly during fetch', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async {
          // Simulate some delay
          await Future.delayed(const Duration(milliseconds: 100));
          return [testPlant1];
        });

        // Act
        final fetchFuture = plantProvider.fetchPlants();
        
        // Assert loading state
        expect(plantProvider.isLoading, isTrue);
        expect(plantProvider.error, isNull);

        await fetchFuture;

        expect(plantProvider.isLoading, isFalse);
      });

      test('uses custom limit and offset parameters', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1]);

        // Act
        await plantProvider.fetchPlants(limit: 10, offset: 5);

        // Assert
        verify(mockPlantService.getUserPlants(limit: 10, offset: 5)).called(1);
      });
    });

    group('getPlantById', () {
      test('returns plant when found', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1, testPlant2]);
        await plantProvider.fetchPlants();

        // Act
        final result = plantProvider.getPlantById(1);

        // Assert
        expect(result, isNotNull);
        expect(result!.id, equals(1));
      });

      test('returns null when plant not found', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1, testPlant2]);
        await plantProvider.fetchPlants();

        // Act
        final result = plantProvider.getPlantById(999);

        // Assert
        expect(result, isNull);
      });
    });

    group('createPlant', () {
      test('successfully creates plant and adds to list', () async {
        // Arrange
        final newPlant = testPlant1.copyWith(id: 3, nickname: 'New Plant');
        when(mockPlantService.createPlant(
          speciesId: anyNamed('speciesId'),
          nickname: anyNamed('nickname'),
          imageUrl: anyNamed('imageUrl'),
          waterIntervalDaysOverride: anyNamed('waterIntervalDaysOverride'),
          fertilizerIntervalDaysOverride: anyNamed('fertilizerIntervalDaysOverride'),
          locationId: anyNamed('locationId'),
          notes: anyNamed('notes'),
        )).thenAnswer((_) async => newPlant);

        // Act
        final result = await plantProvider.createPlant(
          speciesId: 456,
          nickname: 'New Plant',
        );

        // Assert
        expect(result, isNotNull);
        expect(result!.id, equals(3));
        expect(plantProvider.plants, contains(newPlant));
        expect(plantProvider.error, isNull);

        verify(mockPlantService.createPlant(
          speciesId: 456,
          nickname: 'New Plant',
        )).called(1);
      });

      test('handles error during creation', () async {
        // Arrange
        when(mockPlantService.createPlant(
          speciesId: anyNamed('speciesId'),
          nickname: anyNamed('nickname'),
          imageUrl: anyNamed('imageUrl'),
          waterIntervalDaysOverride: anyNamed('waterIntervalDaysOverride'),
          fertilizerIntervalDaysOverride: anyNamed('fertilizerIntervalDaysOverride'),
          locationId: anyNamed('locationId'),
          notes: anyNamed('notes'),
        )).thenThrow(Exception('Creation failed'));

        // Act
        final result = await plantProvider.createPlant(speciesId: 456);

        // Assert
        expect(result, isNull);
        expect(plantProvider.error, contains('Creation failed'));
        expect(plantProvider.plants, isEmpty);
      });
    });

    group('updatePlant', () {
      test('successfully updates plant in list', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1, testPlant2]);
        await plantProvider.fetchPlants();

        final updatedPlant = testPlant1.copyWith(nickname: 'Updated Plant');
        when(mockPlantService.updatePlant(
          plantId: anyNamed('plantId'),
          speciesId: anyNamed('speciesId'),
          nickname: anyNamed('nickname'),
          imageUrl: anyNamed('imageUrl'),
          waterIntervalDaysOverride: anyNamed('waterIntervalDaysOverride'),
          fertilizerIntervalDaysOverride: anyNamed('fertilizerIntervalDaysOverride'),
          locationId: anyNamed('locationId'),
          notes: anyNamed('notes'),
        )).thenAnswer((_) async => updatedPlant);

        // Act
        final result = await plantProvider.updatePlant(
          plantId: 1,
          nickname: 'Updated Plant',
        );

        // Assert
        expect(result, isNotNull);
        expect(result!.nickname, equals('Updated Plant'));
        expect(plantProvider.plants.first.nickname, equals('Updated Plant'));
        expect(plantProvider.error, isNull);

        verify(mockPlantService.updatePlant(
          plantId: 1,
          nickname: 'Updated Plant',
        )).called(1);
      });

      test('handles error during update', () async {
        // Arrange
        when(mockPlantService.updatePlant(
          plantId: anyNamed('plantId'),
          speciesId: anyNamed('speciesId'),
          nickname: anyNamed('nickname'),
          imageUrl: anyNamed('imageUrl'),
          waterIntervalDaysOverride: anyNamed('waterIntervalDaysOverride'),
          fertilizerIntervalDaysOverride: anyNamed('fertilizerIntervalDaysOverride'),
          locationId: anyNamed('locationId'),
          notes: anyNamed('notes'),
        )).thenThrow(Exception('Update failed'));

        // Act
        final result = await plantProvider.updatePlant(plantId: 1);

        // Assert
        expect(result, isNull);
        expect(plantProvider.error, contains('Update failed'));
      });
    });

    group('deletePlant', () {
      test('successfully deletes plant from list', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1, testPlant2]);
        await plantProvider.fetchPlants();

        when(mockPlantService.deletePlant(any)).thenAnswer((_) async => {});

        // Act
        final result = await plantProvider.deletePlant(1);

        // Assert
        expect(result, isTrue);
        expect(plantProvider.plants, hasLength(1));
        expect(plantProvider.plants.first.id, equals(2));
        expect(plantProvider.error, isNull);

        verify(mockPlantService.deletePlant(1)).called(1);
      });

      test('handles error during deletion', () async {
        // Arrange
        when(mockPlantService.deletePlant(any)).thenThrow(Exception('Deletion failed'));

        // Act
        final result = await plantProvider.deletePlant(1);

        // Assert
        expect(result, isFalse);
        expect(plantProvider.error, contains('Deletion failed'));
      });
    });

    group('clearError', () {
      test('clears error state', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenThrow(Exception('Test error'));
        await plantProvider.fetchPlants();
        expect(plantProvider.error, isNotNull);

        // Act
        plantProvider.clearError();

        // Assert
        expect(plantProvider.error, isNull);
      });
    });

    group('refreshPlants', () {
      test('calls fetchPlants', () async {
        // Arrange
        when(mockPlantService.getUserPlants(limit: anyNamed('limit'), offset: anyNamed('offset')))
            .thenAnswer((_) async => [testPlant1]);

        // Act
        await plantProvider.refreshPlants();

        // Assert
        verify(mockPlantService.getUserPlants(limit: 20, offset: 0)).called(1);
        expect(plantProvider.plants, hasLength(1));
      });
    });
  });
}
