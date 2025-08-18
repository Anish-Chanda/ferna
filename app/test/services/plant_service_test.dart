import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:ferna/services/plant_service.dart';
import 'package:ferna/services/http_client.dart';
import 'package:mockito/mockito.dart';
import '../mocks.mocks.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  group('PlantService', () {
    late MockDio mockDio;
    late PlantService plantService;

    const testPlantJson = {
      'id': 1,
      'user_id': 123,
      'species_id': 456,
      'nickname': 'My Plant',
      'image_url': 'https://example.com/image.jpg',
      'watering_frequency_days': 7,
      'last_watered_at': '2025-07-01T10:30:00.000Z',
      'note': 'Living room plant',
      'created_at': '2025-06-01T10:30:00.000Z',
      'updated_at': '2025-07-01T10:30:00.000Z',
    };

    setUp(() async {
      // Create a temporary directory for cookie storage
      final tempDir = Directory.systemTemp.createTempSync();
      final cookieDirPath = '${tempDir.path}/cookies_test/';

      // Initialize HttpClient, passing the temp folder so path_provider isn't needed
      await HttpClient.instance.init(
        baseUrl: 'http://example.com',
        testCookieDir: cookieDirPath,
      );

      // Create and assign the mocked Dio
      mockDio = MockDio();
      HttpClient.instance.dio = mockDio;

      // Grab the singleton PlantService
      plantService = PlantService.instance;
    });

    group('getUserPlants', () {
      test('returns list of plants on successful response', () async {
        // Arrange
        final fakeResponse = Response(
          data: [testPlantJson],
          statusCode: 200,
          requestOptions: RequestOptions(path: '/api/plants'),
        );
        when(
          mockDio.get(
            '/api/plants',
            queryParameters: anyNamed('queryParameters'),
          ),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.getUserPlants();

        // Assert
        expect(result, hasLength(1));
        expect(result.first.id, equals(1));
        expect(result.first.nickname, equals('My Plant'));

        verify(
          mockDio.get(
            '/api/plants',
            queryParameters: {'limit': 20, 'offset': 0},
          ),
        ).called(1);
      });

      test('returns empty list when no plants exist', () async {
        // Arrange
        final fakeResponse = Response(
          data: [],
          statusCode: 200,
          requestOptions: RequestOptions(path: '/api/plants'),
        );
        when(
          mockDio.get(
            '/api/plants',
            queryParameters: anyNamed('queryParameters'),
          ),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.getUserPlants();

        // Assert
        expect(result, isEmpty);
      });

      test('uses custom limit and offset parameters', () async {
        // Arrange
        final fakeResponse = Response(
          data: [testPlantJson],
          statusCode: 200,
          requestOptions: RequestOptions(path: '/api/plants'),
        );
        when(
          mockDio.get(
            '/api/plants',
            queryParameters: anyNamed('queryParameters'),
          ),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        await plantService.getUserPlants(limit: 10, offset: 5);

        // Assert
        verify(
          mockDio.get(
            '/api/plants',
            queryParameters: {'limit': 10, 'offset': 5},
          ),
        ).called(1);
      });

      test('throws exception on non-200 status code', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Server Error',
          statusCode: 500,
          requestOptions: RequestOptions(path: '/api/plants'),
          statusMessage: 'Internal Server Error',
        );
        when(
          mockDio.get(
            '/api/plants',
            queryParameters: anyNamed('queryParameters'),
          ),
        ).thenAnswer((_) async => fakeResponse);

        // Act & Assert
        expect(
          () async => await plantService.getUserPlants(),
          throwsA(isA<Exception>()),
        );
      });
    });

    group('getPlant', () {
      test('returns plant on successful response', () async {
        // Arrange
        final fakeResponse = Response(
          data: testPlantJson,
          statusCode: 200,
          requestOptions: RequestOptions(path: '/api/plants/1'),
        );
        when(mockDio.get('/api/plants/1')).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.getPlant(1);

        // Assert
        expect(result, isNotNull);
        expect(result!.id, equals(1));
        expect(result.nickname, equals('My Plant'));
      });

      test('returns null on 404 response', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Not Found',
          statusCode: 404,
          requestOptions: RequestOptions(path: '/api/plants/999'),
        );
        when(mockDio.get('/api/plants/999')).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.getPlant(999);

        // Assert
        expect(result, isNull);
      });

      test('throws exception on other error status codes', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Server Error',
          statusCode: 500,
          requestOptions: RequestOptions(path: '/api/plants/1'),
          statusMessage: 'Internal Server Error',
        );
        when(mockDio.get('/api/plants/1')).thenAnswer((_) async => fakeResponse);

        // Act & Assert
        expect(
          () async => await plantService.getPlant(1),
          throwsA(isA<Exception>()),
        );
      });
    });

    group('createPlant', () {
      test('creates plant successfully with all parameters', () async {
        // Arrange
        final fakeResponse = Response(
          data: testPlantJson,
          statusCode: 201,
          requestOptions: RequestOptions(path: '/api/plants'),
        );
        when(
          mockDio.post('/api/plants', data: anyNamed('data')),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.createPlant(
          speciesId: 456,
          nickname: 'My Plant',
          imageUrl: 'https://example.com/image.jpg',
          wateringFrequencyDays: 7,
          lastWateredAt: DateTime.parse('2025-07-01T10:30:00.000Z'),
          note: 'Living room plant',
        );

        // Assert
        expect(result.id, equals(1));
        expect(result.nickname, equals('My Plant'));

        verify(
          mockDio.post(
            '/api/plants',
            data: {
              'species_id': 456,
              'nickname': 'My Plant',
              'image_url': 'https://example.com/image.jpg',
              'watering_frequency_days': 7,
              'last_watered_at': '2025-07-01T10:30:00.000Z',
              'note': 'Living room plant',
            },
          ),
        ).called(1);
      });

      test('creates plant with only required parameter', () async {
        // Arrange
        final fakeResponse = Response(
          data: testPlantJson,
          statusCode: 201,
          requestOptions: RequestOptions(path: '/api/plants'),
        );
        when(
          mockDio.post('/api/plants', data: anyNamed('data')),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        await plantService.createPlant(speciesId: 456);

        // Assert
        verify(
          mockDio.post(
            '/api/plants',
            data: {'species_id': 456},
          ),
        ).called(1);
      });

      test('throws exception on non-201 status code', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Bad Request',
          statusCode: 400,
          requestOptions: RequestOptions(path: '/api/plants'),
          statusMessage: 'Bad Request',
        );
        when(
          mockDio.post('/api/plants', data: anyNamed('data')),
        ).thenAnswer((_) async => fakeResponse);

        // Act & Assert
        expect(
          () async => await plantService.createPlant(speciesId: 456),
          throwsA(isA<Exception>()),
        );
      });
    });

    group('updatePlant', () {
      test('updates plant successfully', () async {
        // Arrange
        final updatedPlantJson = Map<String, dynamic>.from(testPlantJson);
        updatedPlantJson['nickname'] = 'Updated Plant';
        
        final fakeResponse = Response(
          data: updatedPlantJson,
          statusCode: 200,
          requestOptions: RequestOptions(path: '/api/plants/1'),
        );
        when(
          mockDio.patch('/api/plants/1', data: anyNamed('data')),
        ).thenAnswer((_) async => fakeResponse);

        // Act
        final result = await plantService.updatePlant(
          plantId: 1,
          nickname: 'Updated Plant',
        );

        // Assert
        expect(result.id, equals(1));
        expect(result.nickname, equals('Updated Plant'));

        verify(
          mockDio.patch(
            '/api/plants/1',
            data: {'nickname': 'Updated Plant'},
          ),
        ).called(1);
      });

      test('throws exception on non-200 status code', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Not Found',
          statusCode: 404,
          requestOptions: RequestOptions(path: '/api/plants/999'),
          statusMessage: 'Not Found',
        );
        when(
          mockDio.patch('/api/plants/999', data: anyNamed('data')),
        ).thenAnswer((_) async => fakeResponse);

        // Act & Assert
        expect(
          () async => await plantService.updatePlant(
            plantId: 999,
            nickname: 'Updated Plant',
          ),
          throwsA(isA<Exception>()),
        );
      });
    });

    group('deletePlant', () {
      test('deletes plant successfully', () async {
        // Arrange
        final fakeResponse = Response(
          data: null,
          statusCode: 204,
          requestOptions: RequestOptions(path: '/api/plants/1'),
        );
        when(mockDio.delete('/api/plants/1')).thenAnswer((_) async => fakeResponse);

        // Act
        await plantService.deletePlant(1);

        // Assert
        verify(mockDio.delete('/api/plants/1')).called(1);
      });

      test('throws exception on non-204 status code', () async {
        // Arrange
        final fakeResponse = Response(
          data: 'Not Found',
          statusCode: 404,
          requestOptions: RequestOptions(path: '/api/plants/999'),
          statusMessage: 'Not Found',
        );
        when(mockDio.delete('/api/plants/999')).thenAnswer((_) async => fakeResponse);

        // Act & Assert
        expect(
          () async => await plantService.deletePlant(999),
          throwsA(isA<Exception>()),
        );
      });
    });
  });
}
