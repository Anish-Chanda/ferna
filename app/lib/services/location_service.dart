import 'package:ferna/services/http_client.dart';
import '../models/location.dart';

class LocationService {
  LocationService._();

  static final LocationService instance = LocationService._();

  // Fetch all locations for the current user
  Future<List<Location>> getUserLocations() async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get('/api/locations');

    if (response.statusCode == 200) {
      final List<dynamic> locationsJson = response.data as List<dynamic>;
      return locationsJson.map((json) => Location.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to fetch locations: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Get a specific location by ID
  Future<Location?> getLocation(int locationId) async {
    final dio = HttpClient.instance.dio;
    
    try {
      final response = await dio.get('/api/locations/$locationId');

      if (response.statusCode == 200) {
        return Location.fromJson(response.data);
      } else if (response.statusCode == 404) {
        return null; // Location not found
      } else {
        throw Exception(
          'Failed to fetch location: HTTP ${response.statusCode}: ${response.statusMessage}',
        );
      }
    } catch (e) {
      throw Exception('Failed to fetch location: $e');
    }
  }

  // Create a new location
  Future<Location> createLocation({
    required String name,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final locationData = <String, dynamic>{
      'name': name,
    };

    final response = await dio.post('/api/locations', data: locationData);

    if (response.statusCode == 201) {
      return Location.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to create location: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Update an existing location
  Future<Location> updateLocation({
    required int locationId,
    required String name,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final updateData = <String, dynamic>{
      'name': name,
    };

    final response = await dio.patch('/api/locations/$locationId', data: updateData);

    if (response.statusCode == 200) {
      return Location.fromJson(response.data);
    } else {
      throw Exception(
        'Failed to update location: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Delete a location
  Future<void> deleteLocation(int locationId) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.delete('/api/locations/$locationId');

    if (response.statusCode != 204) {
      throw Exception(
        'Failed to delete location: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }
}
