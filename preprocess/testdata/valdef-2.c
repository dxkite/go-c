#define A 10 + __LINE__ + B
#define B 20 + A + C
#define C 30 + B

int main() {
    printf("%d\n", C);
}