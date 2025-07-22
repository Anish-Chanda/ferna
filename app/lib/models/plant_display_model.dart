import '../models/plant.dart';

/// A display model that extends Plant with additional computed properties for the UI
class PlantDisplayModel {
  final Plant plant;
  final String? speciesCommonName;
  final String? location; // For now, we'll use the note field as location until we add a proper location field

  PlantDisplayModel({
    required this.plant,
    this.speciesCommonName,
    this.location,
  });

  /// The display name - uses nickname if available, otherwise species common name
  String get displayName => plant.nickname ?? speciesCommonName ?? 'Unknown Plant';

  /// The location text - uses provided location, or shows 'No location' if null
  String get locationText => location ?? 'No location';

  /// Whether this plant has a custom nickname
  bool get hasCustomNickname => plant.nickname != null && plant.nickname!.isNotEmpty;

  /// Calculate days since last watered
  int get daysSinceWatered {
    if (plant.lastWateredAt == null) return 0;
    return DateTime.now().difference(plant.lastWateredAt!).inDays;
  }

  /// Calculate days overdue for watering (negative means not overdue)
  int get daysOverdue => daysSinceWatered - plant.wateringFrequencyDays;

  /// Whether the plant needs watering
  bool get needsWatering => daysOverdue > 0;

  /// Next watering date
  DateTime? get nextWateringDate {
    if (plant.lastWateredAt == null) return null;
    return plant.lastWateredAt!.add(Duration(days: plant.wateringFrequencyDays));
  }
}
