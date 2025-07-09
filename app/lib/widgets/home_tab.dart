import 'package:flutter/material.dart';
import 'package:table_calendar/table_calendar.dart';

class HomeTab extends StatefulWidget {
  const HomeTab({super.key});

  @override
  State<HomeTab> createState() => _HomeTabState();
}

class _HomeTabState extends State<HomeTab> with TickerProviderStateMixin {
  late AnimationController _animationController;
  bool _isExpanded = false;
  DateTime _selectedDate = DateTime.now();
  DateTime _focusedDay = DateTime.now();
  CalendarFormat _calendarFormat = CalendarFormat.week;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 300),
      vsync: this,
    );
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  void _toggleCalendar() {
    setState(() {
      _isExpanded = !_isExpanded;
      _calendarFormat = _isExpanded ? CalendarFormat.month : CalendarFormat.week;
      if (_isExpanded) {
        _animationController.forward();
      } else {
        _animationController.reverse();
      }
    });
  }

  String _getMonthName(DateTime date) {
    const months = [
      'January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'
    ];
    return '${months[date.month - 1]} ${date.year}';
  }

  String _getDateString(DateTime date) {
    const weekdays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
    const months = [
      'January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'
    ];
    return "It's ${weekdays[date.weekday - 1]}, ${months[date.month - 1]} ${date.day}, ${date.year}";
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: Column(
        children: [
          // Calendar Section
          Container(
            decoration: BoxDecoration(
              color: theme.cardTheme.color,
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.1),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: SafeArea(
              child: Column(
                children: [
                  // Header
                  Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          _getMonthName(_focusedDay),
                          style: theme.textTheme.headlineMedium,
                        ),
                        GestureDetector(
                          onTap: () {
                            setState(() {
                              _selectedDate = DateTime.now();
                              _focusedDay = DateTime.now();
                              _calendarFormat = CalendarFormat.week;
                            });
                          },
                          child: Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 6,
                            ),
                            decoration: BoxDecoration(
                              color: theme.chipTheme.backgroundColor,
                              borderRadius: BorderRadius.circular(16),
                            ),
                            child: Text(
                              'Today',
                              style: theme.chipTheme.labelStyle,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  // Calendar
                  Column(
                    children: [
                      GestureDetector(
                        onTap: _toggleCalendar,
                        onVerticalDragUpdate: (details) {
                          // Use vertical drag for more reliable detection
                          if (details.delta.dy > 2 && !_isExpanded) {
                            _toggleCalendar();
                          } else if (details.delta.dy < -2 && _isExpanded) {
                            _toggleCalendar();
                          }
                        },
                        behavior: HitTestBehavior.translucent, // Allow gestures to pass through
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 16.0),
                          child: TableCalendar<dynamic>(
                            firstDay: DateTime.utc(2020, 1, 1),
                            lastDay: DateTime.utc(2030, 12, 31),
                            focusedDay: _focusedDay,
                            calendarFormat: _calendarFormat,
                            availableCalendarFormats: const {
                              CalendarFormat.week: 'Week',
                              CalendarFormat.month: 'Month',
                            },
                            selectedDayPredicate: (day) {
                              return isSameDay(_selectedDate, day);
                            },
                            onDaySelected: (selectedDay, focusedDay) {
                              setState(() {
                                _selectedDate = selectedDay;
                                _focusedDay = focusedDay;
                              });
                            },
                            onPageChanged: (focusedDay) {
                              setState(() {
                                _focusedDay = focusedDay;
                              });
                            },
                            onFormatChanged: (format) {
                              setState(() {
                                _calendarFormat = format;
                                _isExpanded = format == CalendarFormat.month;
                              });
                            },
                            calendarStyle: CalendarStyle(
                              outsideDaysVisible: true,
                              weekendTextStyle: theme.textTheme.bodyMedium!.copyWith(color: theme.colorScheme.onSurface),
                              holidayTextStyle: theme.textTheme.bodyMedium!.copyWith(color: theme.colorScheme.onSurface),
                              selectedDecoration: BoxDecoration(
                                color: theme.colorScheme.primary,
                                shape: BoxShape.circle,
                              ),
                              todayDecoration: BoxDecoration(
                                color: theme.colorScheme.primary.withOpacity(0.3),
                                shape: BoxShape.circle,
                              ),
                              defaultTextStyle: theme.textTheme.bodyMedium!.copyWith(
                                fontWeight: FontWeight.w500,
                              ),
                              outsideTextStyle: theme.textTheme.bodyMedium!.copyWith(
                                color: theme.colorScheme.onSurface.withOpacity(0.4),
                                fontWeight: FontWeight.w500,
                              ),
                              markersMaxCount: 1,
                              markerDecoration: const BoxDecoration(
                                color: Colors.orange,
                                shape: BoxShape.circle,
                              ),
                              markerMargin: const EdgeInsets.only(top: 5),
                              markerSize: 6,
                            ),
                            headerStyle: const HeaderStyle(
                              formatButtonVisible: false,
                              titleCentered: false,
                              headerPadding: EdgeInsets.zero,
                              titleTextStyle: TextStyle(fontSize: 0), // Hide title
                              leftChevronVisible: false,
                              rightChevronVisible: false,
                            ),
                            daysOfWeekStyle: DaysOfWeekStyle(
                              weekdayStyle: theme.textTheme.labelMedium!.copyWith(
                                color: theme.colorScheme.onSurface.withOpacity(0.6),
                              ),
                              weekendStyle: theme.textTheme.labelMedium!.copyWith(
                                color: theme.colorScheme.onSurface.withOpacity(0.6),
                              ),
                            ),
                            eventLoader: (day) {
                              // Add events for certain days (like 7, 9, 11 of current month)
                              if (day.month == _focusedDay.month && [7, 9, 11].contains(day.day)) {
                                return ['event'];
                              }
                              return [];
                            },
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      // Drag indicator - more prominent and responsive
                      GestureDetector(
                        onTap: _toggleCalendar,
                        onVerticalDragUpdate: (details) {
                          // Very sensitive drag detection for the indicator
                          if (details.delta.dy > 1 && !_isExpanded) {
                            _toggleCalendar();
                          } else if (details.delta.dy < -1 && _isExpanded) {
                            _toggleCalendar();
                          }
                        },
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 16), // Larger tap area
                          child: Container(
                            width: 32,
                            height: 4,
                            decoration: BoxDecoration(
                              color: theme.dividerTheme.color,
                              borderRadius: BorderRadius.circular(2),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            ),
          ),
          // Content Section
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  Text(
                    _getDateString(_selectedDate),
                    style: theme.textTheme.bodyLarge,
                  ),
                  const SizedBox(height: 32),
                  // Celebration card
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(24),
                      child: Column(
                        children: [
                          // Plant illustration
                          Container(
                            width: 120,
                            height: 120,
                            child: Image.asset(
                              'assets/illustrations/all_tasks_complete.png',
                              fit: BoxFit.contain,
                            ),
                          ),
                          const SizedBox(height: 16),
                          Text(
                            'Great job!',
                            style: theme.textTheme.headlineLarge,
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Everything is taken care of',
                            style: theme.textTheme.bodyLarge,
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
