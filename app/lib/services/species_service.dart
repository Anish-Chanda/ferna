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
      '/api/species',
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
  // Note: This endpoint doesn't exist in the current backend
  // Use searchSpecies() instead to find species
  Future<Species?> getSpecies(int speciesId) async {
    throw UnimplementedError(
      'getSpecies is not implemented in the backend. Use searchSpecies() instead.'
    );
  }
}
