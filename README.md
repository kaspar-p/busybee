# busybee

![Build Workflow](https://github.com/kaspar-p/busybee/actions/workflows/build.yml/badge.svg)

A discord bot for determining who can get ramen with me. Ingests `.ics` files and is cute.

## Usage

Currently, busybee is manually added to guilds. If you want busybee, contact me.

### Commands

- **.enrol <.ics file>**

  The command to log calendar events with the bot. Simply attach the `.ics` file downloaded from a school site (the ACORN timetable for UofT), and enrol!

- **.wyd <@mention>**

  The command to get today's schedule of a specific person. Note that they have to have enrolled with busybee in the past for this to work. Tells you the titles and times of their commitments.

- **.whenfree <@mention> <@mention> <@mention> ...**

  The command to determine when the users mentioned each have an $N$ hour block of free time in common. Produces a table for each $N$, which ranges from 1 to 6 hours.

- **.whobusy**

  The command to tell you who is currently busy. Only works for users in the system.

## Contributions

If you'd like to contribute to busybee's development, get in touch! Or submit a PR.
