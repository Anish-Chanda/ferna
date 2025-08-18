class Plant {
  final int id;
  final int userId;
  final int speciesId;
  final String? nickname;
  final String? imageUrl;
  final int wateringFrequencyDays;
  final DateTime? lastWateredAt;
  final String? note;
  final DateTime createdAt;
  final DateTime updatedAt;

  Plant({
    required this.id,
    required this.userId,
    required this.speciesId,
    this.nickname,
    this.imageUrl,
    required this.wateringFrequencyDays,
    this.lastWateredAt,
    this.note,
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
      wateringFrequencyDays: json['watering_frequency_days'] as int,
      lastWateredAt: json['last_watered_at'] != null
          ? DateTime.parse(json['last_watered_at'] as String)
          : null,
      note: json['note'] as String?,
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
      'watering_frequency_days': wateringFrequencyDays,
      'last_watered_at': lastWateredAt?.toIso8601String(),
      'note': note,
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
    int? wateringFrequencyDays,
    DateTime? lastWateredAt,
    String? note,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Plant(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      speciesId: speciesId ?? this.speciesId,
      nickname: nickname ?? this.nickname,
      imageUrl: imageUrl ?? this.imageUrl,
      wateringFrequencyDays: wateringFrequencyDays ?? this.wateringFrequencyDays,
      lastWateredAt: lastWateredAt ?? this.lastWateredAt,
      note: note ?? this.note,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'Plant(id: $id, userId: $userId, speciesId: $speciesId, nickname: $nickname, imageUrl: $imageUrl, wateringFrequencyDays: $wateringFrequencyDays, lastWateredAt: $lastWateredAt, note: $note, createdAt: $createdAt, updatedAt: $updatedAt)';
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
        other.wateringFrequencyDays == wateringFrequencyDays &&
        other.lastWateredAt == lastWateredAt &&
        other.note == note &&
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
        wateringFrequencyDays.hashCode ^
        lastWateredAt.hashCode ^
        note.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
