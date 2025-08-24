import 'package:ferna/services/http_client.dart';
import '../models/plant_task.dart';

class PlantTaskService {
  PlantTaskService._();

  static final PlantTaskService instance = PlantTaskService._();

  // Get all tasks for a specific plant
  Future<List<PlantTask>> getTasksForPlant(int plantId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get('/api/plants/$plantId/tasks');

    if (response.statusCode == 200) {
      final List<dynamic> tasksJson = response.data as List<dynamic>;
      return tasksJson.map((json) => PlantTask.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to fetch plant tasks: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Get a specific task by ID
  Future<PlantTask?> getTask(int taskId) async {
    final dio = HttpClient.instance.dio;
    
    try {
      final response = await dio.get('/api/tasks/$taskId');

      if (response.statusCode == 200) {
        return PlantTask.fromJson(response.data);
      } else if (response.statusCode == 404) {
        return null; // Task not found
      } else {
        throw Exception(
          'Failed to fetch task: HTTP ${response.statusCode}: ${response.statusMessage}',
        );
      }
    } catch (e) {
      throw Exception('Failed to fetch task: $e');
    }
  }

  // Get overdue tasks for the current user
  Future<List<PlantTask>> getOverdueTasks() async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get('/api/tasks/overdue');

    if (response.statusCode == 200) {
      final List<dynamic> tasksJson = response.data as List<dynamic>;
      return tasksJson.map((json) => PlantTask.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to fetch overdue tasks: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Create a new task for a plant
  Future<PlantTask> createTask({
    required int plantId,
    required String taskType, // 'watering' or 'fertilizer'
    required int intervalDays,
    int toleranceDays = 0,
    DateTime? snoozedUntil,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final taskData = <String, dynamic>{
      'task_type': taskType,
      'interval_days': intervalDays,
      'tolerance_days': toleranceDays,
      if (snoozedUntil != null) 'snoozed_until': snoozedUntil.toIso8601String(),
    };

    final response = await dio.post('/api/plants/$plantId/tasks', data: taskData);

    if (response.statusCode == 201) {
      return PlantTask.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to create task: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Update an existing task
  Future<PlantTask> updateTask({
    required int taskId,
    String? taskType,
    int? intervalDays,
    int? toleranceDays,
    DateTime? snoozedUntil,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final updateData = <String, dynamic>{};
    if (taskType != null) updateData['task_type'] = taskType;
    if (intervalDays != null) updateData['interval_days'] = intervalDays;
    if (toleranceDays != null) updateData['tolerance_days'] = toleranceDays;
    if (snoozedUntil != null) updateData['snoozed_until'] = snoozedUntil.toIso8601String();

    final response = await dio.patch('/api/tasks/$taskId', data: updateData);

    if (response.statusCode == 200) {
      return PlantTask.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to update task: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Delete a task
  Future<void> deleteTask(int taskId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.delete('/api/tasks/$taskId');

    if (response.statusCode != 204) {
      throw Exception(
        'Failed to delete task: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }
}
