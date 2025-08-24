import 'package:ferna/services/http_client.dart';
import '../models/plant.dart';

class PlantService {
  PlantService._();

  static final PlantService instance = PlantService._();

  // Fetch all plants for the current user
  Future<List<Plant>> getUserPlants({int limit = 20, int offset = 0}) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get(
      '/api/plants',
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );

    if (response.statusCode == 200) {
      final List<dynamic> plantsJson = response.data as List<dynamic>;
      return plantsJson.map((json) => Plant.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to fetch plants: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Get a specific plant by ID
  Future<Plant?> getPlant(int plantId) async {
    final dio = HttpClient.instance.dio;
    
    try {
      final response = await dio.get('/api/plants/$plantId');

      if (response.statusCode == 200) {
        return Plant.fromJson(response.data);
      } else if (response.statusCode == 404) {
        return null; // Plant not found
      } else {
        throw Exception(
          'Failed to fetch plant: HTTP ${response.statusCode}: ${response.statusMessage}',
        );
      }
    } catch (e) {
      throw Exception('Failed to fetch plant: $e');
    }
  }

  // Create a new plant
  Future<Plant> createPlant({
    required int speciesId,
    String? nickname,
    String? imageUrl,
    int? waterIntervalDaysOverride,
    int? fertilizerIntervalDaysOverride,
    int? locationId,
    String? notes,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final plantData = <String, dynamic>{
      'species_id': speciesId,
      if (nickname != null) 'nickname': nickname,
      if (imageUrl != null) 'image_url': imageUrl,
      if (waterIntervalDaysOverride != null) 'water_interval_days_override': waterIntervalDaysOverride,
      if (fertilizerIntervalDaysOverride != null) 'fertilizer_interval_days_override': fertilizerIntervalDaysOverride,
      if (locationId != null) 'location_id': locationId,
      if (notes != null) 'notes': notes,
    };

    final response = await dio.post('/api/plants', data: plantData);

    if (response.statusCode == 201) {
      return Plant.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to create plant: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Update an existing plant
  Future<Plant> updatePlant({
    required int plantId,
    int? speciesId,
    String? nickname,
    String? imageUrl,
    int? waterIntervalDaysOverride,
    int? fertilizerIntervalDaysOverride,
    int? locationId,
    String? notes,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final updateData = <String, dynamic>{};
    if (speciesId != null) updateData['species_id'] = speciesId;
    if (nickname != null) updateData['nickname'] = nickname;
    if (imageUrl != null) updateData['image_url'] = imageUrl;
    if (waterIntervalDaysOverride != null) updateData['water_interval_days_override'] = waterIntervalDaysOverride;
    if (fertilizerIntervalDaysOverride != null) updateData['fertilizer_interval_days_override'] = fertilizerIntervalDaysOverride;
    if (locationId != null) updateData['location_id'] = locationId;
    if (notes != null) updateData['notes'] = notes;

    final response = await dio.patch('/api/plants/$plantId', data: updateData);

    if (response.statusCode == 200) {
      return Plant.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to update plant: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Delete a plant
  Future<void> deletePlant(int plantId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.delete('/api/plants/$plantId');

    if (response.statusCode != 204) {
      throw Exception(
        'Failed to delete plant: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }
}
