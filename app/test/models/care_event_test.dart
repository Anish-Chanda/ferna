import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/models/care_event.dart';

void main() {
  group('CareEvent Model', () {
    const testCareEventJson = {
      'id': 1,
      'plant_id': 123,
      'task_id': 456,
      'event_type': 'watering',
      'happened_at': '2025-06-01T10:30:00.000Z',
      'notes': 'Watered thoroughly',
      'created_at': '2025-06-01T10:00:00.000Z',
      'updated_at': '2025-06-01T11:00:00.000Z',
    };

    const testCareEventJsonMinimal = {
      'id': 2,
      'plant_id': 456,
      'event_type': 'fertilizer',
      'happened_at': '2025-06-01T10:30:00.000Z',
      'created_at': '2025-06-01T10:00:00.000Z',
      'updated_at': '2025-06-01T10:00:00.000Z',
    };

    test('creates CareEvent from JSON with all fields', () {
      final careEvent = CareEvent.fromJson(testCareEventJson);

      expect(careEvent.id, equals(1));
      expect(careEvent.plantId, equals(123));
      expect(careEvent.taskId, equals(456));
      expect(careEvent.eventType, equals(EventType.watering));
      expect(careEvent.happenedAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(careEvent.notes, equals('Watered thoroughly'));
      expect(careEvent.createdAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
      expect(careEvent.updatedAt, equals(DateTime.parse('2025-06-01T11:00:00.000Z')));
    });

    test('creates CareEvent from JSON with minimal fields (nullable fields as null)', () {
      final careEvent = CareEvent.fromJson(testCareEventJsonMinimal);

      expect(careEvent.id, equals(2));
      expect(careEvent.plantId, equals(456));
      expect(careEvent.taskId, isNull);
      expect(careEvent.eventType, equals(EventType.fertilizer));
      expect(careEvent.happenedAt, equals(DateTime.parse('2025-06-01T10:30:00.000Z')));
      expect(careEvent.notes, isNull);
      expect(careEvent.createdAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
      expect(careEvent.updatedAt, equals(DateTime.parse('2025-06-01T10:00:00.000Z')));
    });

    test('converts CareEvent to JSON', () {
      final careEvent = CareEvent(
        id: 1,
        plantId: 123,
        taskId: 456,
        eventType: EventType.watering,
        happenedAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        notes: 'Watered thoroughly',
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final json = careEvent.toJson();

      expect(json['id'], equals(1));
      expect(json['plant_id'], equals(123));
      expect(json['task_id'], equals(456));
      expect(json['event_type'], equals('watering'));
      expect(json['happened_at'], equals('2025-06-01T10:30:00.000Z'));
      expect(json['notes'], equals('Watered thoroughly'));
      expect(json['created_at'], equals('2025-06-01T10:00:00.000Z'));
      expect(json['updated_at'], equals('2025-06-01T11:00:00.000Z'));
    });

    test('copyWith creates new instance with updated fields', () {
      final originalEvent = CareEvent(
        id: 1,
        plantId: 123,
        taskId: null,
        eventType: EventType.watering,
        happenedAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        notes: null,
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
      );

      final updatedEvent = originalEvent.copyWith(
        notes: 'Added some fertilizer',
        taskId: 789,
      );

      expect(updatedEvent.id, equals(originalEvent.id));
      expect(updatedEvent.plantId, equals(originalEvent.plantId));
      expect(updatedEvent.taskId, equals(789));
      expect(updatedEvent.eventType, equals(originalEvent.eventType));
      expect(updatedEvent.happenedAt, equals(originalEvent.happenedAt));
      expect(updatedEvent.notes, equals('Added some fertilizer'));
      expect(updatedEvent.createdAt, equals(originalEvent.createdAt));
      expect(updatedEvent.updatedAt, equals(originalEvent.updatedAt));
    });

    test('equality works correctly', () {
      final event1 = CareEvent(
        id: 1,
        plantId: 123,
        taskId: 456,
        eventType: EventType.watering,
        happenedAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        notes: 'Watered thoroughly',
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final event2 = CareEvent(
        id: 1,
        plantId: 123,
        taskId: 456,
        eventType: EventType.watering,
        happenedAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        notes: 'Watered thoroughly',
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T11:00:00.000Z'),
      );

      final event3 = event1.copyWith(notes: 'Different notes');

      expect(event1, equals(event2));
      expect(event1, isNot(equals(event3)));
      expect(event1.hashCode, equals(event2.hashCode));
      expect(event1.hashCode, isNot(equals(event3.hashCode)));
    });

    test('toString returns meaningful representation', () {
      final event = CareEvent(
        id: 1,
        plantId: 123,
        taskId: null,
        eventType: EventType.watering,
        happenedAt: DateTime.parse('2025-06-01T10:30:00.000Z'),
        notes: null,
        createdAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
        updatedAt: DateTime.parse('2025-06-01T10:00:00.000Z'),
      );

      final stringRepresentation = event.toString();
      
      expect(stringRepresentation, contains('CareEvent('));
      expect(stringRepresentation, contains('id: 1'));
      expect(stringRepresentation, contains('plantId: 123'));
      expect(stringRepresentation, contains('eventType: EventType.watering'));
    });

    test('EventType enum conversion works correctly', () {
      expect(EventType.watering.toString(), equals('EventType.watering'));
      expect(EventType.fertilizer.toString(), equals('EventType.fertilizer'));
      expect(EventType.repotting.toString(), equals('EventType.repotting'));
      expect(EventType.pruning.toString(), equals('EventType.pruning'));
      expect(EventType.other.toString(), equals('EventType.other'));
    });
  });
}
