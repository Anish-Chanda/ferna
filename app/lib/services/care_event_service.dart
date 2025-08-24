import 'package:ferna/services/http_client.dart';
import '../models/care_event.dart';

class CareEventService {
  CareEventService._();

  static final CareEventService instance = CareEventService._();

  // Get all care events for a specific plant
  Future<List<CareEvent>> getEventsForPlant(int plantId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get('/api/plants/$plantId/events');

    if (response.statusCode == 200) {
      final List<dynamic> eventsJson = response.data as List<dynamic>;
      return eventsJson.map((json) => CareEvent.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to fetch care events: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Get a specific care event by ID
  Future<CareEvent?> getEvent(int eventId) async {
    final dio = HttpClient.instance.dio;
    
    try {
      final response = await dio.get('/api/events/$eventId');

      if (response.statusCode == 200) {
        return CareEvent.fromJson(response.data);
      } else if (response.statusCode == 404) {
        return null; // Event not found
      } else {
        throw Exception(
          'Failed to fetch care event: HTTP ${response.statusCode}: ${response.statusMessage}',
        );
      }
    } catch (e) {
      throw Exception('Failed to fetch care event: $e');
    }
  }

  // Create a new care event for a plant
  Future<CareEvent> createEvent({
    required int plantId,
    required String eventType, // 'watering', 'fertilizer', 'repotting', 'pruning', 'other'
    int? taskId,
    DateTime? happenedAt,
    String? notes,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final eventData = <String, dynamic>{
      'event_type': eventType,
      if (taskId != null) 'task_id': taskId,
      'happened_at': (happenedAt ?? DateTime.now()).toIso8601String(),
      if (notes != null) 'notes': notes,
    };

    final response = await dio.post('/api/plants/$plantId/events', data: eventData);

    if (response.statusCode == 201) {
      return CareEvent.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to create care event: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Update an existing care event
  Future<CareEvent> updateEvent({
    required int eventId,
    String? eventType,
    int? taskId,
    DateTime? happenedAt,
    String? notes,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final updateData = <String, dynamic>{};
    if (eventType != null) updateData['event_type'] = eventType;
    if (taskId != null) updateData['task_id'] = taskId;
    if (happenedAt != null) updateData['happened_at'] = happenedAt.toIso8601String();
    if (notes != null) updateData['notes'] = notes;

    final response = await dio.patch('/api/events/$eventId', data: updateData);

    if (response.statusCode == 200) {
      return CareEvent.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to update care event: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Delete a care event
  Future<void> deleteEvent(int eventId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.delete('/api/events/$eventId');

    if (response.statusCode != 204) {
      throw Exception(
        'Failed to delete care event: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }
}
