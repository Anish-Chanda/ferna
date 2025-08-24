enum TaskType {
  watering,
  fertilizer,
}

class PlantTask {
  final int id;
  final int plantId;
  final TaskType taskType;
  final DateTime? snoozedUntil;
  final int intervalDays;
  final int toleranceDays;
  final DateTime? nextDueAt;
  final DateTime createdAt;
  final DateTime updatedAt;

  PlantTask({
    required this.id,
    required this.plantId,
    required this.taskType,
    this.snoozedUntil,
    required this.intervalDays,
    required this.toleranceDays,
    this.nextDueAt,
    required this.createdAt,
    required this.updatedAt,
  });

  factory PlantTask.fromJson(Map<String, dynamic> json) {
    return PlantTask(
      id: json['id'] as int,
      plantId: json['plant_id'] as int,
      taskType: _taskTypeFromString(json['task_type'] as String),
      snoozedUntil: json['snoozed_until'] != null
          ? DateTime.parse(json['snoozed_until'] as String)
          : null,
      intervalDays: json['interval_days'] as int,
      toleranceDays: json['tolerance_days'] as int,
      nextDueAt: json['next_due_at'] != null
          ? DateTime.parse(json['next_due_at'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'plant_id': plantId,
      'task_type': _taskTypeToString(taskType),
      'snoozed_until': snoozedUntil?.toIso8601String(),
      'interval_days': intervalDays,
      'tolerance_days': toleranceDays,
      'next_due_at': nextDueAt?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  static TaskType _taskTypeFromString(String value) {
    switch (value) {
      case 'watering':
        return TaskType.watering;
      case 'fertilizer':
        return TaskType.fertilizer;
      default:
        throw ArgumentError('Unknown task type: $value');
    }
  }

  static String _taskTypeToString(TaskType taskType) {
    switch (taskType) {
      case TaskType.watering:
        return 'watering';
      case TaskType.fertilizer:
        return 'fertilizer';
    }
  }

  PlantTask copyWith({
    int? id,
    int? plantId,
    TaskType? taskType,
    DateTime? snoozedUntil,
    int? intervalDays,
    int? toleranceDays,
    DateTime? nextDueAt,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return PlantTask(
      id: id ?? this.id,
      plantId: plantId ?? this.plantId,
      taskType: taskType ?? this.taskType,
      snoozedUntil: snoozedUntil ?? this.snoozedUntil,
      intervalDays: intervalDays ?? this.intervalDays,
      toleranceDays: toleranceDays ?? this.toleranceDays,
      nextDueAt: nextDueAt ?? this.nextDueAt,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'PlantTask(id: $id, plantId: $plantId, taskType: $taskType, snoozedUntil: $snoozedUntil, intervalDays: $intervalDays, toleranceDays: $toleranceDays, nextDueAt: $nextDueAt, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;

    return other is PlantTask &&
        other.id == id &&
        other.plantId == plantId &&
        other.taskType == taskType &&
        other.snoozedUntil == snoozedUntil &&
        other.intervalDays == intervalDays &&
        other.toleranceDays == toleranceDays &&
        other.nextDueAt == nextDueAt &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt;
  }

  @override
  int get hashCode {
    return id.hashCode ^
        plantId.hashCode ^
        taskType.hashCode ^
        snoozedUntil.hashCode ^
        intervalDays.hashCode ^
        toleranceDays.hashCode ^
        nextDueAt.hashCode ^
        createdAt.hashCode ^
        updatedAt.hashCode;
  }
}
