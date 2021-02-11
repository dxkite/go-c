#define A 10 + __LINE__ + B
#define B 20 + A + C
#define C 30 + B

#define D 'D' + E
#define E 'E' + F
#define F 'F'

int main() {
    printf("%d %d\n", C, E);
}