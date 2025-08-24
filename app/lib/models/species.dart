class Species {
  final int id;
  final String commonName;
  final String? scientificName;
  final String lightPreference;
  final int defaultWaterIntervalDays;
  final int defaultFertilizerIntervalDays;
  final String toxicity;
  final String? careNotes;
  final String? careNotesSource;
  final DateTime createdAt;
  final DateTime updatedAt;

  Species({
    required this.id,
    required this.commonName,
    this.scientificName,
    required this.lightPreference,
    required this.defaultWaterIntervalDays,
    required this.defaultFertilizerIntervalDays,
    required this.toxicity,
    this.careNotes,
    this.careNotesSource,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Species.fromJson(Map<String, dynamic> json) {
    return Species(
      id: json['id'] as int,
      commonName: json['common_name'] as String,
      scientificName: json['scientific_name'] as String?,
      lightPreference: json['light_pref'] as String,
      defaultWaterIntervalDays: json['default_water_interval_days'] as int,
      defaultFertilizerIntervalDays: json['default_fertilizer_interval_days'] as int,
      toxicity: json['toxicity'] as String,
      careNotes: json['care_notes'] as String?,
      careNotesSource: json['care_notes_source'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'common_name': commonName,
      'scientific_name': scientificName,
      'light_pref': lightPreference,
      'default_water_interval_days': defaultWaterIntervalDays,
      'default_fertilizer_interval_days': defaultFertilizerIntervalDays,
      'toxicity': toxicity,
      'care_notes': careNotes,
      'care_notes_source': careNotesSource,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  Species copyWith({
    int? id,
    String? commonName,
    String? scientificName,
    String? lightPreference,
    int? defaultWaterIntervalDays,
    int? defaultFertilizerIntervalDays,
    String? toxicity,
    String? careNotes,
    String? careNotesSource,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Species(
      id: id ?? this.id,
      commonName: commonName ?? this.commonName,
      scientificName: scientificName ?? this.scientificName,
      lightPreference: lightPreference ?? this.lightPreference,
      defaultWaterIntervalDays: defaultWaterIntervalDays ?? this.defaultWaterIntervalDays,
      defaultFertilizerIntervalDays: defaultFertilizerIntervalDays ?? this.defaultFertilizerIntervalDays,
      toxicity: toxicity ?? this.toxicity,
      careNotes: careNotes ?? this.careNotes,
      careNotesSource: careNotesSource ?? this.careNotesSource,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'Species(id: $id, commonName: $commonName, scientificName: $scientificName, lightPreference: $lightPreference, defaultWaterIntervalDays: $defaultWaterIntervalDays, defaultFertilizerIntervalDays: $defaultFertilizerIntervalDays, toxicity: $toxicity, careNotes: $careNotes, careNotesSource: $careNotesSource, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;

    return other is Species &&
        other.id == id &&
        other.commonName == commonName &&
        other.scientificName == scientificName &&
        other.lightPreference == lightPreference &&
        other.defaultWaterIntervalDays == defaultWaterIntervalDays &&
        other.defaultFertilizerIntervalDays == defaultFertilizerIntervalDays &&
        other.toxicity == toxicity &&
        other.careNotes == careNotes &&
        other.careNotesSource == careNotesSource &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        commonName.hashCode ^
        scientificName.hashCode ^
        lightPreference.hashCode ^
        defaultWaterIntervalDays.hashCode ^
        defaultFertilizerIntervalDays.hashCode ^
        toxicity.hashCode ^
        careNotes.hashCode ^
        careNotesSource.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
