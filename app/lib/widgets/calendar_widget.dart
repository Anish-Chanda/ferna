import 'package:flutter/material.dart';
import 'package:table_calendar/table_calendar.dart';

class CalendarWidget extends StatefulWidget {
  final DateTime focusedDay;
  final DateTime selectedDate;
  final CalendarFormat calendarFormat;
  final Function(DateTime, DateTime) onDaySelected;
  final Function(DateTime) onPageChanged;
  final Function(CalendarFormat) onFormatChanged;

  const CalendarWidget({
    super.key,
    required this.focusedDay,
    required this.selectedDate,
    required this.calendarFormat,
    required this.onDaySelected,
    required this.onPageChanged,
    required this.onFormatChanged,
  });

  @override
  State<CalendarWidget> createState() => _CalendarWidgetState();
}

class _CalendarWidgetState extends State<CalendarWidget> {
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return TableCalendar<dynamic>(
      firstDay: DateTime.utc(2020, 1, 1),
      lastDay: DateTime.utc(2030, 12, 31),
      focusedDay: widget.focusedDay,
      calendarFormat: widget.calendarFormat,
      availableCalendarFormats: const {
        CalendarFormat.week: 'Week',
        CalendarFormat.month: 'Month',
      },
      selectedDayPredicate: (day) {
        return isSameDay(widget.selectedDate, day);
      },
      onDaySelected: widget.onDaySelected,
      onPageChanged: widget.onPageChanged,
      onFormatChanged: widget.onFormatChanged,
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
        // TODO: events can be added here
        return [];
      },
    );
  }
}
