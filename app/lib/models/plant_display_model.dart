import '../models/plant.dart';
import '../models/species.dart';

class PlantDisplayModel {
  final Plant plant;
  final Species? species;
  final String? location;

  PlantDisplayModel({
    required this.plant,
    this.species,
    this.location,
  });

  String get displayName => plant.nickname ?? species?.commonName ?? 'Unknown Plant';

  String get locationText => location ?? 'No location';

  bool get hasCustomNickname => plant.nickname != null && plant.nickname!.isNotEmpty;

  // Since we no longer have lastWateredAt, we'll need to calculate this differently
  // For now, returning a default value until we implement care event tracking
  int get daysSinceWatered => 0;

  // Calculate based on species default or plant override
  int get wateringFrequencyDays {
    if (plant.waterIntervalDaysOverride != null) {
      return plant.waterIntervalDaysOverride!;
    }
    return species?.defaultWaterIntervalDays ?? 7;
  }

  int get daysOverdue => daysSinceWatered - wateringFrequencyDays;

  bool get needsWatering => daysOverdue > 0;

  // This will need to be updated when we implement care event tracking
  DateTime? get nextWateringDate => null;
}
