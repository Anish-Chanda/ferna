import '../models/plant.dart';

class PlantDisplayModel {
  final Plant plant;
  final String? speciesCommonName;
  final String? location;

  PlantDisplayModel({
    required this.plant,
    this.speciesCommonName,
    this.location,
  });

  String get displayName => plant.nickname ?? speciesCommonName ?? 'Unknown Plant';

  String get locationText => location ?? 'No location';

  bool get hasCustomNickname => plant.nickname != null && plant.nickname!.isNotEmpty;

  int get daysSinceWatered {
    if (plant.lastWateredAt == null) return 0;
    return DateTime.now().difference(plant.lastWateredAt!).inDays;
  }

  int get daysOverdue => daysSinceWatered - plant.wateringFrequencyDays;

  bool get needsWatering => daysOverdue > 0;

  DateTime? get nextWateringDate {
    if (plant.lastWateredAt == null) return null;
    return plant.lastWateredAt!.add(Duration(days: plant.wateringFrequencyDays));
  }
}
