#define A() A+10+B+__COUNTER__
#define B(x,y) A() + 10 + __LINE__ + #x + #y + C(1,2,3)
#define C(...) # __VA_ARGS__
#define D(a,b, C, ...) a+b+C(__VA_ARGS__, __LINE__)+B(__LINE__, bb)
#define E 123 ## abc
#define S_INT # 10
A();
B(19, 29);
B(,20);
C(a, b, c);
D(1,2,3, A,B, C(a, b, c), B(19, 29) ,D,E, __LINE__);
E
S_INT

#define hash_hash # ## #
#define mkstr(a) # a
#define in_between(a) mkstr(a)
#define join(c, d) in_between(c hash_hash d)
char p[] = join(x, y); // equivalent to
// char p[] = "x ## y";