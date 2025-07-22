class Species {
  final int id;
  final String commonName;
  final String? scientificName;
  final int defaultWateringFrequency;
  final DateTime createdAt;
  final DateTime updatedAt;

  Species({
    required this.id,
    required this.commonName,
    this.scientificName,
    required this.defaultWateringFrequency,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Species.fromJson(Map<String, dynamic> json) {
    return Species(
      id: json['id'] as int,
      commonName: json['common_name'] as String,
      scientificName: json['scientific_name'] as String?,
      defaultWateringFrequency: json['default_watering_frequency_days'] as int,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'common_name': commonName,
      'scientific_name': scientificName,
      'default_watering_frequency_days': defaultWateringFrequency,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  Species copyWith({
    int? id,
    String? commonName,
    String? scientificName,
    int? defaultWateringFrequency,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Species(
      id: id ?? this.id,
      commonName: commonName ?? this.commonName,
      scientificName: scientificName ?? this.scientificName,
      defaultWateringFrequency: defaultWateringFrequency ?? this.defaultWateringFrequency,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'Species(id: $id, commonName: $commonName, scientificName: $scientificName, defaultWateringFrequency: $defaultWateringFrequency, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;

    return other is Species &&
        other.id == id &&
        other.commonName == commonName &&
        other.scientificName == scientificName &&
        other.defaultWateringFrequency == defaultWateringFrequency &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        commonName.hashCode ^
        scientificName.hashCode ^
        defaultWateringFrequency.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
