enum EventType {
  watering,
  fertilizer,
  repotting,
  pruning,
  other,
}

class CareEvent {
  final int id;
  final int plantId;
  final int? taskId;
  final EventType eventType;
  final DateTime happenedAt;
  final String? notes;
  final DateTime createdAt;
  final DateTime updatedAt;

  CareEvent({
    required this.id,
    required this.plantId,
    this.taskId,
    required this.eventType,
    required this.happenedAt,
    this.notes,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CareEvent.fromJson(Map<String, dynamic> json) {
    return CareEvent(
      id: json['id'] as int,
      plantId: json['plant_id'] as int,
      taskId: json['task_id'] as int?,
      eventType: _eventTypeFromString(json['event_type'] as String),
      happenedAt: DateTime.parse(json['happened_at'] as String),
      notes: json['notes'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'plant_id': plantId,
      'task_id': taskId,
      'event_type': _eventTypeToString(eventType),
      'happened_at': happenedAt.toIso8601String(),
      'notes': notes,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  static EventType _eventTypeFromString(String value) {
    switch (value) {
      case 'watering':
        return EventType.watering;
      case 'fertilizer':
        return EventType.fertilizer;
      case 'repotting':
        return EventType.repotting;
      case 'pruning':
        return EventType.pruning;
      case 'other':
        return EventType.other;
      default:
        throw ArgumentError('Unknown event type: $value');
    }
  }

  static String _eventTypeToString(EventType eventType) {
    switch (eventType) {
      case EventType.watering:
        return 'watering';
      case EventType.fertilizer:
        return 'fertilizer';
      case EventType.repotting:
        return 'repotting';
      case EventType.pruning:
        return 'pruning';
      case EventType.other:
        return 'other';
    }
  }

  CareEvent copyWith({
    int? id,
    int? plantId,
    int? taskId,
    EventType? eventType,
    DateTime? happenedAt,
    String? notes,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return CareEvent(
      id: id ?? this.id,
      plantId: plantId ?? this.plantId,
      taskId: taskId ?? this.taskId,
      eventType: eventType ?? this.eventType,
      happenedAt: happenedAt ?? this.happenedAt,
      notes: notes ?? this.notes,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'CareEvent(id: $id, plantId: $plantId, taskId: $taskId, eventType: $eventType, happenedAt: $happenedAt, notes: $notes, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;

    return other is CareEvent &&
        other.id == id &&
        other.plantId == plantId &&
        other.taskId == taskId &&
        other.eventType == eventType &&
        other.happenedAt == happenedAt &&
        other.notes == notes &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        plantId.hashCode ^
        taskId.hashCode ^
        eventType.hashCode ^
        happenedAt.hashCode ^
        notes.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
