import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/species.dart';

void main() {
  group('Species Model', () {
    const testSpeciesJson = {
      'id': 1,
      'common_name': 'Monstera Deliciosa',
      'scientific_name': 'Monstera deliciosa',
      'light_pref': 'bright_indirect',
      'default_water_interval_days': 7,
      'default_fertilizer_interval_days': 30,
      'toxicity': 'non_toxic',
      'care_notes': 'Keep soil moist but not soggy',
      'care_notes_source': 'Horticulture guide',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    const testSpeciesJsonMinimal = {
      'id': 2,
      'common_name': 'Snake Plant',
      'light_pref': 'low',
      'default_water_interval_days': 14,
      'default_fertilizer_interval_days': 60,
      'toxicity': 'toxic_to_pets',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    test('creates Species from JSON with all fields', () {
      final species = Species.fromJson(testSpeciesJson);

      expect(species.id, equals(1));
      expect(species.commonName, equals('Monstera Deliciosa'));
      expect(species.scientificName, equals('Monstera deliciosa'));
      expect(species.lightPreference, equals('bright_indirect'));
      expect(species.defaultWaterIntervalDays, equals(7));
      expect(species.defaultFertilizerIntervalDays, equals(30));
      expect(species.toxicity, equals('non_toxic'));
      expect(species.careNotes, equals('Keep soil moist but not soggy'));
      expect(species.careNotesSource, equals('Horticulture guide'));
      expect(species.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(species.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('creates Species from JSON with minimal fields (scientific name as null)', () {
      final species = Species.fromJson(testSpeciesJsonMinimal);

      expect(species.id, equals(2));
      expect(species.commonName, equals('Snake Plant'));
      expect(species.scientificName, isNull);
      expect(species.lightPreference, equals('low'));
      expect(species.defaultWaterIntervalDays, equals(14));
      expect(species.defaultFertilizerIntervalDays, equals(60));
      expect(species.toxicity, equals('toxic_to_pets'));
      expect(species.careNotes, isNull);
      expect(species.careNotesSource, isNull);
      expect(species.createdAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(species.updatedAt, equals(DateTime.parse('2025-07-01T10:30:00.000Z')));
    });

    test('converts Species to JSON', () {
      final species = Species(
        id: 1,
        commonName: 'Monstera Deliciosa',
        scientificName: 'Monstera deliciosa',
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        careNotes: 'Keep soil moist but not soggy',
        careNotesSource: 'Horticulture guide',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final json = species.toJson();

      expect(json['id'], equals(1));
      expect(json['common_name'], equals('Monstera Deliciosa'));
      expect(json['scientific_name'], equals('Monstera deliciosa'));
      expect(json['light_pref'], equals('bright_indirect'));
      expect(json['default_water_interval_days'], equals(7));
      expect(json['default_fertilizer_interval_days'], equals(30));
      expect(json['toxicity'], equals('non_toxic'));
      expect(json['care_notes'], equals('Keep soil moist but not soggy'));
      expect(json['care_notes_source'], equals('Horticulture guide'));
      expect(json['created_at'], equals('2025-06-01T10:30:00.000Z'));
      expect(json['updated_at'], equals('2025-07-01T10:30:00.000Z'));
    });

    test('copyWith creates new instance with updated fields', () {
      final originalSpecies = Species(
        id: 1,
        commonName: 'Original Name',
        scientificName: null,
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final updatedSpecies = originalSpecies.copyWith(
        commonName: 'Updated Name',
        scientificName: 'Updated scientificus',
        defaultWaterIntervalDays: 10,
      );

      expect(updatedSpecies.id, equals(originalSpecies.id));
      expect(updatedSpecies.commonName, equals('Updated Name'));
      expect(updatedSpecies.scientificName, equals('Updated scientificus'));
      expect(updatedSpecies.defaultWaterIntervalDays, equals(10));
      expect(updatedSpecies.lightPreference, equals(originalSpecies.lightPreference));
      expect(updatedSpecies.createdAt, equals(originalSpecies.createdAt));
      expect(updatedSpecies.updatedAt, equals(originalSpecies.updatedAt));
    });

    test('equality works correctly', () {
      final species1 = Species(
        id: 1,
        commonName: 'Monstera Deliciosa',
        scientificName: 'Monstera deliciosa',
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final species2 = Species(
        id: 1,
        commonName: 'Monstera Deliciosa',
        scientificName: 'Monstera deliciosa',
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final species3 = species1.copyWith(commonName: 'Different Name');

      expect(species1, equals(species2));
      expect(species1, isNot(equals(species3)));
      expect(species1.hashCode, equals(species2.hashCode));
      expect(species1.hashCode, isNot(equals(species3.hashCode)));
    });

    test('toString returns meaningful representation', () {
      final species = Species(
        id: 1,
        commonName: 'Monstera Deliciosa',
        scientificName: 'Monstera deliciosa',
        lightPreference: 'bright_indirect',
        defaultWaterIntervalDays: 7,
        defaultFertilizerIntervalDays: 30,
        toxicity: 'non_toxic',
        createdAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        updatedAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
      );

      final stringRepresentation = species.toString();
      
      expect(stringRepresentation, contains('Species('));
      expect(stringRepresentation, contains('id: 1'));
      expect(stringRepresentation, contains('commonName: Monstera Deliciosa'));
      expect(stringRepresentation, contains('scientificName: Monstera deliciosa'));
    });
  });
}
