import 'package:ferna/services/http_client.dart';
import '../models/species.dart';

class SpeciesService {
  SpeciesService._();

  static final SpeciesService instance = SpeciesService._();

  // Search for species by query
  Future<List<Species>> searchSpecies({
    required String query,
    int limit = 20,
    int offset = 0,
  }) async {
    final dio = HttpClient.instance.dio;
    
    final response = await dio.get(
      '/api/plants/species',
      queryParameters: {
        'query': query,
        'limit': limit,
        'offset': offset,
      },
    );

    if (response.statusCode == 200) {
      final List<dynamic> speciesJson = response.data as List<dynamic>;
      return speciesJson.map((json) => Species.fromJson(json)).toList();
    } else {
      throw Exception(
        'Failed to search species: HTTP ${response.statusCode}: ${response.statusMessage}',
      );
    }
  }

  // Get a specific species by ID
  Future<Species?> getSpecies(int speciesId) async {
    final dio = HttpClient.instance.dio;
    
    try {
      final response = await dio.get('/api/plants/species/$speciesId');

      if (response.statusCode == 200) {
        return Species.fromJson(response.data);
      } else if (response.statusCode == 404) {
        return null; // Species not found
      } else {
        throw Exception(
          'Failed to fetch species: HTTP ${response.statusCode}: ${response.statusMessage}',
        );
      }
    } catch (e) {
      throw Exception('Failed to fetch species: $e');
    }
  }
}
