import 'package:flutter/material.dart';
import '../models/plant.dart';
import '../services/plant_service.dart';

class PlantProvider with ChangeNotifier {
  final PlantService _plantService;
  
  // Constructor that allows injecting a custom PlantService (useful for testing)
  PlantProvider({PlantService? plantService})
      : _plantService = plantService ?? PlantService.instance;
  
  List<Plant> _plants = [];
  List<Plant> get plants => _plants;
  
  bool _isLoading = false;
  bool get isLoading => _isLoading;
  
  String? _error;
  String? get error => _error;

  // Fetch all plants for the current user
  Future<void> fetchPlants({int limit = 20, int offset = 0}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _plants = await _plantService.getUserPlants(limit: limit, offset: offset);
      _error = null;
    } catch (e) {
      _error = e.toString();
      _plants = [];
    }

    _isLoading = false;
    notifyListeners();
  }

  // Get a specific plant by ID
  Plant? getPlantById(int id) {
    try {
      return _plants.firstWhere((plant) => plant.id == id);
    } catch (e) {
      return null;
    }
  }

  // Create a new plant
  Future<Plant?> createPlant({
    required int speciesId,
    String? nickname,
    String? imageUrl,
    int? waterIntervalDaysOverride,
    int? fertilizerIntervalDaysOverride,
    int? locationId,
    String? notes,
  }) async {
    try {
      final newPlant = await _plantService.createPlant(
        speciesId: speciesId,
        nickname: nickname,
        imageUrl: imageUrl,
        waterIntervalDaysOverride: waterIntervalDaysOverride,
        fertilizerIntervalDaysOverride: fertilizerIntervalDaysOverride,
        locationId: locationId,
        notes: notes,
      );
      
      _plants.add(newPlant);
      notifyListeners();
      return newPlant;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  // Update an existing plant
  Future<Plant?> updatePlant({
    required int plantId,
    int? speciesId,
    String? nickname,
    String? imageUrl,
    int? waterIntervalDaysOverride,
    int? fertilizerIntervalDaysOverride,
    int? locationId,
    String? notes,
  }) async {
    try {
      final updatedPlant = await _plantService.updatePlant(
        plantId: plantId,
        speciesId: speciesId,
        nickname: nickname,
        imageUrl: imageUrl,
        waterIntervalDaysOverride: waterIntervalDaysOverride,
        fertilizerIntervalDaysOverride: fertilizerIntervalDaysOverride,
        locationId: locationId,
        notes: notes,
      );
      
      final index = _plants.indexWhere((plant) => plant.id == plantId);
      if (index != -1) {
        _plants[index] = updatedPlant;
        notifyListeners();
      }
      return updatedPlant;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  // Delete a plant
  Future<bool> deletePlant(int plantId) async {
    try {
      await _plantService.deletePlant(plantId);
      _plants.removeWhere((plant) => plant.id == plantId);
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  // Clear error message
  void clearError() {
    _error = null;
    notifyListeners();
  }

  // Refresh plants list
  Future<void> refreshPlants() async {
    await fetchPlants();
  }
}