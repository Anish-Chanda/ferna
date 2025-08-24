class Plant {
  final int id;
  final int userId;
  final int speciesId;
  final String? nickname;
  final String? imageUrl;
  final String? notes;
  final int? waterIntervalDaysOverride;
  final int? fertilizerIntervalDaysOverride;
  final int? locationId;
  final DateTime createdAt;
  final DateTime updatedAt;

  Plant({
    required this.id,
    required this.userId,
    required this.speciesId,
    this.nickname,
    this.imageUrl,
    this.notes,
    this.waterIntervalDaysOverride,
    this.fertilizerIntervalDaysOverride,
    this.locationId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Plant.fromJson(Map<String, dynamic> json) {
    return Plant(
      id: json['id'] as int,
      userId: json['user_id'] as int,
      speciesId: json['species_id'] as int,
      nickname: json['nickname'] as String?,
      imageUrl: json['image_url'] as String?,
      notes: json['notes'] as String?,
      waterIntervalDaysOverride: json['water_interval_days_override'] as int?,
      fertilizerIntervalDaysOverride: json['fertilizer_interval_days_override'] as int?,
      locationId: json['location_id'] as int?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'species_id': speciesId,
      'nickname': nickname,
      'image_url': imageUrl,
      'notes': notes,
      'water_interval_days_override': waterIntervalDaysOverride,
      'fertilizer_interval_days_override': fertilizerIntervalDaysOverride,
      'location_id': locationId,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  Plant copyWith({
    int? id,
    int? userId,
    int? speciesId,
    String? nickname,
    String? imageUrl,
    String? notes,
    int? waterIntervalDaysOverride,
    int? fertilizerIntervalDaysOverride,
    int? locationId,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Plant(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      speciesId: speciesId ?? this.speciesId,
      nickname: nickname ?? this.nickname,
      imageUrl: imageUrl ?? this.imageUrl,
      notes: notes ?? this.notes,
      waterIntervalDaysOverride: waterIntervalDaysOverride ?? this.waterIntervalDaysOverride,
      fertilizerIntervalDaysOverride: fertilizerIntervalDaysOverride ?? this.fertilizerIntervalDaysOverride,
      locationId: locationId ?? this.locationId,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'Plant(id: $id, userId: $userId, speciesId: $speciesId, nickname: $nickname, imageUrl: $imageUrl, notes: $notes, waterIntervalDaysOverride: $waterIntervalDaysOverride, fertilizerIntervalDaysOverride: $fertilizerIntervalDaysOverride, locationId: $locationId, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;

    return other is Plant &&
        other.id == id &&
        other.userId == userId &&
        other.speciesId == speciesId &&
        other.nickname == nickname &&
        other.imageUrl == imageUrl &&
        other.notes == notes &&
        other.waterIntervalDaysOverride == waterIntervalDaysOverride &&
        other.fertilizerIntervalDaysOverride == fertilizerIntervalDaysOverride &&
        other.locationId == locationId &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        userId.hashCode ^
        speciesId.hashCode ^
        nickname.hashCode ^
        imageUrl.hashCode ^
        notes.hashCode ^
        waterIntervalDaysOverride.hashCode ^
        fertilizerIntervalDaysOverride.hashCode ^
        locationId.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
