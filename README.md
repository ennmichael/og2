### Notes

The initial build time might be a bit long because I'm using SQLite and it needs
to be built with gcc. (https://github.com/mattn/go-sqlite3#installation)

To save time, I only implemented up to level three, instead of five levels. Hopefully
I will be forgiven for this - all that's needed to get to level five is to extend
the `factory.properties` method to include the other two levels as well, but this
is a bit laborious.

I'm serializing and storing the game state as JSON - this should be good enough
for our current requirements. In production, we might want to have proper columns
for fields of the state, should we need them.

All games are always kept in memory. Again, this might cause problems in production,
if there are too many active games and memory runs out. The solution is either always
serializing and loading games, or sharding the simulation.

In the interest of time and simplicity, some errors are simply paniced on instead
of being propagated and handled or more elegantly displayed to the user.

I use SQLite for simplicity, which might not be a good choice for production.
Something like Postgres would be better, or if we anticipate really big scale
maybe some NoSQL solution or sharded Postgres. However if we only expect very
low loads (e.g. this is a side project or a demo), SQLite serves as a good choice
even when going to production.

Normally I use TDD, but this time there are no tests due to the tight time constraint.
