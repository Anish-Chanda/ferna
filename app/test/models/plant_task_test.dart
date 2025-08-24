import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/plant_task.dart';

void main() {
  group('PlantTask Model', () {
    const testPlantTaskJson = {
      'id': 1,
      'plant_id': 123,
      'task_type': 'watering',
      'snoozed_until': '2025-06-02T10:30:00.000Z',
      'interval_days': 7,
      'tolerance_days': 2,
      'next_due_at': '2025-06-08T10:30:00.000Z',
      'created_at': '2025-06-01T10:00:00.000Z',
      'updated_at': '2025-06-01T11:00:00.000Z',
    };

    const testPlantTaskJsonMinimal = {
      'id': 2,
      'plant_id': 456,
      'task_type': 'fertilizer',
      'interval_days': 30,
      'tolerance_days': 5,
      'created_at': '2025-06-01T10:00:00.000Z',
      'updated_at': '2025-06-01T10:00:00.000Z',
    };

    test('creates PlantTask from JSON with all fields', () {
      final plantTask = PlantTask.fromJson(testPlantTaskJson);

      expect(plantTask.id, equals(1));
      expect(plantTask.plantId, equals(123));
      expect(plantTask.taskType, equals(TaskType.watering));
      expect(plantTask.snoozedUntil, equals(DateTime.parse('2025-06-02T10:30:00.000Z')));
      expect(plantTask.intervalDays, equals(7));
      expect(plantTask.toleranceDays, equals(2));
      expect(plantTask.nextDueAt, equals(DateTime.parse('2025-06-08T10:30:00.000Z')));
      expect(plantTask.createdAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
      expect(plantTask.updatedAt, equals(DateTime.parse('2025-06-01T11:00:00.000Z')));
    });

    test('creates PlantTask from JSON with minimal fields (nullable fields as null)', () {
      final plantTask = PlantTask.fromJson(testPlantTaskJsonMinimal);

      expect(plantTask.id, equals(2));
      expect(plantTask.plantId, equals(456));
      expect(plantTask.taskType, equals(TaskType.fertilizer));
      expect(plantTask.snoozedUntil, isNull);
      expect(plantTask.intervalDays, equals(30));
      expect(plantTask.toleranceDays, equals(5));
      expect(plantTask.nextDueAt, isNull);
      expect(plantTask.createdAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
      expect(plantTask.updatedAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
    });

    test('converts PlantTask to JSON', () {
      final plantTask = PlantTask(
        id: 1,
        plantId: 123,
        taskType: TaskType.watering,
        snoozedUntil: DateTime.parse('2025-06-02T10:30:00.000Z'),
        intervalDays: 7,
        toleranceDays: 2,
        nextDueAt: DateTime.parse('2025-06-08T10:30:00.000Z'),
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final json = plantTask.toJson();

      expect(json['id'], equals(1));
      expect(json['plant_id'], equals(123));
      expect(json['task_type'], equals('watering'));
      expect(json['snoozed_until'], equals('2025-06-02T10:30:00.000Z'));
      expect(json['interval_days'], equals(7));
      expect(json['tolerance_days'], equals(2));
      expect(json['next_due_at'], equals('2025-06-08T10:30:00.000Z'));
      expect(json['created_at'], equals('2025-06-01T10:00:00.000Z'));
      expect(json['updated_at'], equals('2025-06-01T11:00:00.000Z'));
    });

    test('copyWith creates new instance with updated fields', () {
      final originalTask = PlantTask(
        id: 1,
        plantId: 123,
        taskType: TaskType.watering,
        snoozedUntil: null,
        intervalDays: 7,
        toleranceDays: 2,
        nextDueAt: null,
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
      );

      final updatedTask = originalTask.copyWith(
        snoozedUntil: DateTime.parse('2025-06-02T10:30:00.000Z'),
        nextDueAt: DateTime.parse('2025-06-08T10:30:00.000Z'),
      );

      expect(updatedTask.id, equals(originalTask.id));
      expect(updatedTask.plantId, equals(originalTask.plantId));
      expect(updatedTask.taskType, equals(originalTask.taskType));
      expect(updatedTask.snoozedUntil, equals(DateTime.parse('2025-06-02T10:30:00.000Z')));
      expect(updatedTask.intervalDays, equals(originalTask.intervalDays));
      expect(updatedTask.toleranceDays, equals(originalTask.toleranceDays));
      expect(updatedTask.nextDueAt, equals(DateTime.parse('2025-06-08T10:30:00.000Z')));
      expect(updatedTask.createdAt, equals(originalTask.createdAt));
      expect(updatedTask.updatedAt, equals(originalTask.updatedAt));
    });

    test('equality works correctly', () {
      final task1 = PlantTask(
        id: 1,
        plantId: 123,
        taskType: TaskType.watering,
        snoozedUntil: DateTime.parse('2025-06-02T10:30:00.000Z'),
        intervalDays: 7,
        toleranceDays: 2,
        nextDueAt: DateTime.parse('2025-06-08T10:30:00.000Z'),
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final task2 = PlantTask(
        id: 1,
        plantId: 123,
        taskType: TaskType.watering,
        snoozedUntil: DateTime.parse('2025-06-02T10:30:00.000Z'),
        intervalDays: 7,
        toleranceDays: 2,
        nextDueAt: DateTime.parse('2025-06-08T10:30:00.000Z'),
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final task3 = task1.copyWith(intervalDays: 14);

      expect(task1, equals(task2));
      expect(task1, isNot(equals(task3)));
      expect(task1.hashCode, equals(task2.hashCode));
      expect(task1.hashCode, isNot(equals(task3.hashCode)));
    });

    test('toString returns meaningful representation', () {
      final task = PlantTask(
        id: 1,
        plantId: 123,
        taskType: TaskType.watering,
        snoozedUntil: null,
        intervalDays: 7,
        toleranceDays: 2,
        nextDueAt: null,
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
      );

      final stringRepresentation = task.toString();
      
      expect(stringRepresentation, contains('PlantTask('));
      expect(stringRepresentation, contains('id: 1'));
      expect(stringRepresentation, contains('plantId: 123'));
      expect(stringRepresentation, contains('taskType: TaskType.watering'));
    });

    test('TaskType enum conversion works correctly', () {
      expect(TaskType.watering.toString(), equals('TaskType.watering'));
      expect(TaskType.fertilizer.toString(), equals('TaskType.fertilizer'));
    });
  });
}
