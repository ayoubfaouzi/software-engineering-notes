# A Tour of Rust

## Rustup and Cargo

- `rustup` is the recommended Rust installation and toolchain manager. It installs Rust across supported platforms and simplifies upgrades through commands such as `rustup update`.
- A standard installation provides three main commands:
  - `cargo` manages projects, dependencies, compilation, and program execution.
  - `rustc` is the compiler, normally invoked indirectly by Cargo but available for direct use.
  - `rustdoc` generates HTML documentation from specially formatted source comments, also typically invoked through Cargo.

| Command                 | What it does                                                                |
| ----------------------- | --------------------------------------------------------------------------- |
| `cargo new <name>`      | Create a new project in a new directory                                     |
| `cargo init`            | Create a project in the current directory                                   |
| `cargo build`           | Compile the project (debug mode), fetching dependencies if needed           |
| `cargo build --release` | Compile with optimizations                                                  |
| `cargo run`             | Build and run the binary                                                    |
| `cargo run -- <args>`   | Build and run, passing arguments to your program                            |
| `cargo check`           | Check that the code compiles without producing a binary (faster than build) |
| `cargo test`            | Build and run tests                                                         |
| `cargo clean`           | Delete the `target` directory (build artifacts)                             |
| `cargo doc --open`      | Build documentation for your project and its dependencies, then open it     |
| `cargo fmt`             | Format your code                                                            |
| `cargo clippy`          | Run the linter for common mistakes and style issues                         |
| `cargo install <crate>` | Install a binary crate globally (for example, `cargo-outdated`)             |

## Rust Functions

```rust
fn gcd(mut n: u64, mut m: u64) -> u64 {
    while m != 0 {
        if m < n {
            let t = m;
            m = n;
            n = t;
        }
        m = m % n;
    }
    n
}
```

- The `fn` keyword (pronounced “fun”) introduces a function.
- **Four-space** indentation is standard Rust style.
- By **default**, once a variable is **initialized**, its value **can’t be changed**, but placing the `mut` keyword (pronounced “mute,” short for mutable) before the parameters n and m allows our function body to assign to them.
- Unlike C and C‍++, Rust does **not require** parentheses around the **conditional** expressions, but it does require curly braces around the statements they control.
- A `let` statement declares a local variable, like `t` in our function. We don’t need to write out t’s type, as long as Rust can infer it from how the variable is used.
  - If we wanted to spell out t’s type, we could write: `let t: u64 = m;`
- Rust has a `return` statement, but the `gcd` function doesn’t need one. If a function body ends with an expression that is not followed by a semicolon, that’s the function’s return value.
  - It’s typical in Rust to use this form to establish the function’s value when control “falls off the end” of the function
  - Use `return` statements only for explicit early returns from the **midst** of a function.

## Writing and Running Unit Tests

```rust
#[test]
fn test_gcd() {
    assert_eq!(gcd(14, 15), 1);
    assert_eq!(gcd(2 * 3 * 54321, 5 * 7 * 54321), 54321);
    assert_eq!(gcd(0, 0), 0);
}
```

- The `#[test]` atop the definition marks `test_gcd` as a test function, to be skipped in **normal compilations**, but included and called automatically if we running the `cargo test` command.
- The `#[test]` marker is an example of an **attribute**. Attributes are an open-ended system for marking functions and other declarations with extra information, like attributes in C‍++ and C#, or annotations in Java.
- The ! character marks these as a **macro calls**, not **function calls**. During compilation, macro calls are “expanded” to Rust code.
- Unlike C and C‍++, in which assertions can be **skipped**, Rust **always** checks assertions regardless of how the program was compiled 👍.
- There are also `debug_assert!` and `debug_assert_eq!` macros, whose assertions are skipped when the program is compiled for speed.

## Handling Command-Line Arguments

```rust
use std::str::FromStr;
use std::env;

// Our main function doesn’t return a value, so we can simply omit the -> and return type.
fn main() {
    let mut numbers = Vec::new();

    for arg in env::args().skip(1) {
        // Here we call `u64::from_str` to attempt to parse our command-line argument arg as an unsigned 64-bit integer.
        // Rather than a method we’re invoking on a particular u64 value, `u64::from_str` is a function associated with
        // the u64 type, akin to a static method in C‍++ or Java.
        numbers.push(u64::from_str(&arg).expect("error parsing argument"));
    }

    if numbers.len() == 0 {
        eprintln!("Usage: gcd NUMBER ..."); // write our error message to the standard error output stream.
        std::process::exit(1);
    }

    let mut d = numbers[0];
    // ownership of the vector should remain with numbers, we are merely borrowing its elements for the loop.
    for m in &numbers[1..] {
        d = gcd(d, *m);
    }
    // since numbers owns the vector, Rust automatically frees it when numbers goes out of scope at the end of main.

    println!("The greatest common divisor of {numbers:?} is {d}");

    // Alternatively, we could have our main function return a std::process::ExitCode or a Result,
    // in which case the return value would determine the exit status. In the case of an error Result,
    // Rust would also print the debug form of the error to stderr.
}
```

- The first `use` declaration brings the **standard library trait** `FromStr` into scope.
- A **trait** is a collection of **function signatures** that types can implement, like an interface in Java or C#.
- Any type that implements the `FromStr` trait has a `from_str` method that tries to parse a value of that type from a string.
- We declare a mutable local variable `numbers` and initialize it to an empty vector. `Vec` is Rust’s growable vector type, analogous to C‍++’s `std::vector`.

> 👉 The type of numbers is `Vec<u64>`, a vector of `u64` values, but as before, we don’t need to write that out. Rust will **infer** it for us, in part because what we push onto the vector are `u64` values, but also because we pass the vector’s elements to `gcd`, which accepts only `u64` values.

- The `std::env` module’s `args` function returns an **iterator** over the command-line arguments. That is, its return type implements the `Iterator` trait, which provides methods for looping over a sequence of values.
- The `from_str` function doesn’t return a u64 directly, but rather a `Result` value that indicates whether the parse succeeded or failed. A `Result` value is one of two variants:
  - A value written `Ok(v)`, indicating that the parse succeeded and v is the value produced.
  - A value written `Err(e)`, indicating that the parse failed and e is an error value explaining why.

## Serving Pages to the Web

- A Rust package, whether a **library** or an **executable**, is called a **crate**; Cargo and crates.io both derive their names from this term.
- To add a dependency, you can either do a `cargo add package-name@a.b.c` or edit the TOML file directly.

| Command                 | What it does                                                                |
| ----------------------- | --------------------------------------------------------------------------- |
| `cargo fetch`           | Download dependencies without compiling                                     |
| `cargo add <crate>`     | Add a dependency to `Cargo.toml`                                            |
| `cargo remove <crate>`  | Remove a dependency from `Cargo.toml`                                       |
| `cargo update`          | Update dependencies to the newest versions allowed by `Cargo.toml`          |
| `cargo update <crate>`  | Update a single dependency                                                  |
| `cargo install <crate>` | Install a binary crate globally (for example, `cargo-outdated`)             |
| `cargo tree`            | Show the dependency tree                                                    |

Crates can have **optional** features: parts of the interface or implementation that not all users need, and which are therefore excluded from the build by default: `$ cargo add serde@1.0.228 --features derive`.

```rust
// When we write use actix_web::{...}, each of the names listed inside the curly brackets
// becomes directly usable in our code; instead of having to spell out the full name
// actix_web::HttpResponse each time we use it, we can simply refer to it as HttpResponse.
use actix_web::{App, HttpResponse, HttpServer, web};

// Actix is written using asynchronous code to support serving thousands of connections
// at a time without spawning thousands of system threads. 
#[actix_web::main]
async fn main() {
    let server = HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(get_index)) //
        // Since there’s no semicolon at the end of the closure’s body,
        // the App is the closure’s return value.
    });

    println!("Serving on http://localhost:3000...");
    server
        .bind("127.0.0.1:3000")
        .expect("error binding server to address")
        .run()
        .await
        .expect("error running server");
}

async fn get_index() -> HttpResponse {
    HttpResponse::Ok()
        // We call its content_type and body methods to fill in the details of the response;
        // each call returns the HttpResponse it was applied to, with the modifications made.
        .content_type("text/html")
        .body(
            r#"
                <title>GCD Calculator</title>
                <form action="/gcd" method="post">
                <input type="text" name="n"/>
                <input type="text" name="m"/>
                <button type="submit">Compute GCD</button>
                </form>
            "#
        )
}
```

- The argument we pass to `HttpServer::new` is the Rust **closure** expression `|| { App::new() ... }`.
- A closure is a value that can be called as if it were a function.
  - This closure takes no arguments, but if it did, their names would appear between the `||` vertical bars.
  - The `{ ... }` is the body of the closure.
- When we start our server, `Actix` starts a pool of threads to handle incoming requests. Each thread calls our closure to get a fresh copy of the `App` value that tells it how to route and handle requests.
- Since the response text contains a lot of **double quotes**, we write it using the Rust “raw string” syntax.
  - the letter `r`, zero or more hash marks, a double quote, and then the contents of the string, terminated by another double quote followed by the same number of hash marks.
  - Any character may occur within a raw string **without being escaped**, including double quotes; in fact, no escape sequences like \" are recognized.
- Next, let’s define a Rust structure type that represents the values we expect from our form:
```rust
// Placing a #[derive(Deserialize)] attribute above a type definition tells the serde crate
// to examine the type when the program is compiled and automatically generate code to parse
// a value of this type from data in the format that HTML forms use for POST requests.
#[derive(Deserialize)]
struct GcdParameters {
    n: u64,
    m: u64, // The comma after m: u64 is optional, since it’s the last field. 
}
```
- With this definition in place, we can write our handler function quite easily:
```rust
async fn post_gcd(form: web::Form<GcdParameters>) -> HttpResponse {
    if form.n == 0 || form.m == 0 {
        return HttpResponse::BadRequest()
            .content_type("text/html")
            .body("Computing the GCD with zero is boring.");
    }

    // The format! macro is just like the println! macro, except that instead
    // of writing the text to the standard output, it returns it as a string.
    // Since the values we want to print are not simple variable names,
    // we must use empty braces {} to mark the place where we want to insert
    // them, then pass the values as extra arguments. 
    let response = format!(
        "The greatest common divisor of the numbers {} and {} \
            is <b>{}</b>\n",
        form.n,
        form.m,
        gcd(form.n, form.m),
    );

    HttpResponse::Ok()
        .content_type("text/html")
        .body(response)
}
```

## Concurrency

Rust’s ownership and borrowing rules prevent both memory errors and data races: mutex-protected data can be accessed only while holding the lock, which is released automatically; shared read-only data cannot be accidentally modified; and transferring ownership between threads prevents the sender from accessing the transferred data again. In C and C++, these guarantees depend more heavily on programmer discipline.

### What the Mandelbrot Set Actually Is

The Mandelbrot set is defined as the set of complex numbers `c` for which `z` does not fly out to infinity.
```rust
use num::Complex;

fn complex_square_add_loop(c: Complex<f64>) {
    let mut z = Complex { re: 0.0, im: 0.0 };
    loop {
        z = z * z + c;
    }
}
```

Complex is a Rust structure type (or struct), defined like this:
```rust
// Complex is a generic structure: you can read the <T> after the type name as “for any type T.” 
struct Complex<T> {
    /// Real portion of the complex number
    re: T,

    /// Imaginary portion of the complex number
    im: T,
}
```

- The infinite loop takes a while to run, but there are two tricks for the impatient:
  - First, if we give up on running the loop forever and just try some limited number of iterations, it turns out that we still get a decent approximation of the set. How many iterations we need depends on how precisely we want to plot the boundary. 
  - Second, it’s been shown that, if z ever once leaves the circle of **radius 2** centered at the origin, it will **definitely fly** infinitely far away from the origin eventually. 

So here’s the final version of our loop, and the heart of our program:

```rust
use num::Complex;

/// Tries to determine if `c` is in the Mandelbrot set, using at most `limit`
/// iterations to decide.
///
/// If `c` is not a member, this returns `Some(i)`, where `i` is the number of
/// iterations it took for `c` to leave the circle of radius 2 centered on the
/// origin. If `c` seems to be a member (more precisely, if we reached the
/// iteration limit without being able to prove that `c` is not a member),
/// this returns `None`.
fn escape_time(c: Complex<f64>, limit: usize) -> Option<usize> {
    let mut z = Complex { re: 0.0, im: 0.0 };
    for i in 0..limit {
        // Instead of computing a square root, we just compare the squared
        // distance with 4.0, which is faster.
        if z.norm_sqr() > 4.0 {
            return Some(i);
        }
        z = z * z + c;
    }

    None
}
```

The function’s return value is an `Option<usize>`. Rust’s standard library defines the Option type as follows:
```rust
// Option is a generic type: you can use Option<T> to represent an optional value of any type T you like.
enum Option<T> {
    None,
    Some(T),
}
```

###  Parsing Pair Command-Line Arguments

```rust
use std::str::FromStr;

/// Parses the string `s` as a coordinate pair, like `"400x600"` or `"1.0,0.5"`.
///
/// Specifically, `s` should have the form <left><sep><right>, where <sep> is
/// the character given by the `separator` argument, and <left> and <right> are
/// both strings that can be parsed by `T::from_str`. `separator` must be an
/// ASCII character.
///
/// If `s` has the proper form, this returns `Some<(x, y)>`. If it doesn't parse
/// correctly, this returns `None`.
fn parse_pair<T: FromStr>(s: &str, separator: char) -> Option<(T, T)> {
    match s.find(separator) {
        None => None,
        Some(index) => {
            match (T::from_str(&s[..index]), T::from_str(&s[index + 1..])) {
                (Ok(l), Ok(r)) => Some((l, r)), // This pattern matches only if both
                                                // Results are Ok variants, indicating
                                                // that both parses succeeded.
                _ => None, // The wildcard pattern _ matches anything and ignores its
                           // value. If we reach this point, then parse_pair has failed,
                           // so we evaluate to None, again providing the return value of the function.
            }
        }
    }
}

#[test]
fn test_parse_pair() {
    // Rust will often be able to infer type parameters for you,
    // and you won’t need to write them out as we did in the test code.
    assert_eq!(parse_pair::<i32>("",        ','), None);
    assert_eq!(parse_pair::<i32>("10,",     ','), None);
    assert_eq!(parse_pair::<i32>(",10",     ','), None);
    assert_eq!(parse_pair::<i32>("10,20",   ','), Some((10, 20)));
    assert_eq!(parse_pair::<i32>("10,20xy", ','), None);
    assert_eq!(parse_pair::<f64>("0.5x",    'x'), None);
    assert_eq!(parse_pair::<f64>("0.5x1.5", 'x'), Some((0.5, 1.5)));
}
```

- The definition of `parse_pair` is a generic function, you can read the clause `<T: FromStr>` aloud as, *“For any type T that implements the FromStr trait...”*.
- Rust programmer would call `T` a **type parameter** of `parse_pair`.
- Our return type is` Option<(T, T)>`: either `None` or a value `Some((v1, v2))`, where (`v1, v2)` is a **tuple** of two values, both of type `T`.
- The `parse_pair` function doesn’t use an **explicit return** statement, here, that's the whole `match` expression, so whichever arm runs produces the result.

Now that we have parse_pair, it’s easy to write a function to parse a pair of floating-point coordinates and return them as a Complex<f64> value:
```rust
/// Parses a pair of floating-point numbers separated by a comma as a
/// complex number.
fn parse_complex(s: &str) -> Option<Complex<f64>> {
    match parse_pair(s, ',') {
        Some((re, im)) => Some(Complex { re, im }),
        None => None,
    }
}

#[test]
fn test_parse_complex() {
    assert_eq!(
        parse_complex("1.25,-0.0625"),
        Some(Complex { re: 1.25, im: -0.0625 }),
    );
    assert_eq!(parse_complex(",-0.0625"), None);
}
```

- If you were reading closely, you may have noticed that we used a shorthand notation to build the `Complex` value.
- It’s common to initialize a struct’s fields with variables of the **same name**, so rather than forcing you to write `Complex { re: re, im: im }`, Rust lets you simply write `Complex { re, im }`. This is modeled on similar notations in JavaScript and Haskell.

### Mapping from Pixels to Complex Numbers

```rust
/// Computes the point on the complex plane that corresponds to a given
/// pixel in the output image.
///
/// `bounds` is a pair giving the width and height of the image in pixels.
/// `pixel` is a (column, row) pair indicating a particular pixel in that image.
/// The `upper_left` and `lower_right` parameters are points on the complex
/// plane designating the area our image covers.
fn pixel_to_point(
    bounds: (usize, usize),
    pixel: (usize, usize),
    upper_left: Complex<f64>,
    lower_right: Complex<f64>,
) -> Complex<f64> {
    let (width, height) = (
        lower_right.re - upper_left.re,
        upper_left.im - lower_right.im
    );
    Complex {
        re: upper_left.re + pixel.0 as f64 * width  / bounds.0 as f64,
        im: upper_left.im - pixel.1 as f64 * height / bounds.1 as f64,
        // Why subtraction here? pixel.1 increases as we go down,
        // but the imaginary component increases as we go up.
    }
}

#[test]
fn test_pixel_to_point() {
    assert_eq!(
        pixel_to_point(
            (100, 200),
            (25, 175),
            Complex { re: -1.0, im: 1.0 },
            Complex { re: 1.0, im: -1.0 },
        ),
        Complex { re: -0.5, im: -0.75 },
    );
}
```

- Expressions with this form refer to tuple fields: `pixel.0`.
  - This refers to the **first** field of the tuple pixel: `pixel.0 as f64`.
  - This is Rust’s syntax for a **type conversion**: this converts `pixel.0` to an `f64` value. Unlike C and C‍++, Rust generally refuses to convert between numeric types implicitly !

### Plotting the Set

```rust
/// Renders a rectangle of the Mandelbrot set into a buffer of pixels.
///
/// The `bounds` argument gives the width and height of the buffer `pixels`,
/// which holds one grayscale pixel per byte. The `upper_left` and `lower_right`
/// arguments specify points on the complex plane corresponding to the upper-
/// left and lower-right corners of the pixel buffer.
fn render(
    pixels: &mut [u8],
    bounds: (usize, usize),
    upper_left: Complex<f64>,
    lower_right: Complex<f64>,
) {
    assert!(pixels.len() == bounds.0 * bounds.1);

    for row in 0..bounds.1 {
        for column in 0..bounds.0 {
            let point =
                pixel_to_point(bounds, (column, row), upper_left, lower_right);
            pixels[row * bounds.0 + column] =
                match escape_time(point, 255) {
                    // If escape_time says that point belongs to the set, render 
                    // colors the corresponding pixel black (0). Otherwise, render
                    // assigns darker colors to the numbers that took longer to
                    // escape the circle.
                    None => 0,
                    Some(count) => 255 - count as u8,
                };
        }
    }
}
```


- The type of the first argument, pixels, is `&mut [u8]`. 
  - In English, that’s a mutable reference, `&mut`, to a **slice** of unsigned 8-bit integers, [u8].
  - Think of a slice as an array or vector; pixels will refer to a **slab of contiguous memory**, typically millions of `u8` bytes.
  - We need a mutable reference in order to write to the buffer, since Rust references are **read-only** by default 🤷.

### Writing Image Files

```rust
use image::{ExtendedColorType, ImageEncoder, ImageError};
use image::codecs::png::PngEncoder;
use std::fs::File;

/// Writes the buffer `pixels`, whose dimensions are given by `bounds`, to the
/// file named `filename`.
fn write_image(
    filename: &str,
    pixels: &[u8],
    bounds: (usize, usize),
) -> Result<(), ImageError> {
    let output = File::create(filename)?;

    let encoder = PngEncoder::new(output);
    encoder.write_image(
        pixels,
        bounds.0 as u32,
        bounds.1 as u32,
        ExtendedColorType::L8,
    )?;

    Ok(())
}
```

- When all goes well, our `write_image` function has no useful value to return; it wrote everything interesting to the file. So its success type is the *unit* type `()`, so called because it has only one value, also written `()`. The unit type is akin to *void* in C and C‍++.
- The return type of `File::create` is `Result<std::fs::File, std::io::Error>`, while that of `encoder.encode` is `Result<(), std::io::Error>`, so both share the same error type, `std::io::Error`. It makes sense for our `write_image` function to do the same. In either case, failure should result in an immediate return, passing along the `std::io::Error` value describing what went wrong.
- One way to handle `File::create`’s result would be to match on its return value, like this:
```rust
let output = match File::create(filename) {
    Ok(f) => f, // On success, let output be the File carried in the Ok value. 
    Err(e) => { // On failure, pass the error along to our own caller.
        return Err(e);
    }
};
```
- This kind of match statement is such a common pattern in Rust that the language provides the ❓ operator as shorthand for the whole thing.
-  So, rather than writing out this logic explicitly in `write_image`, we used the following equivalent and much more legible statement:
```rust
// If File::create fails, the ? operator returns from write_image,
// passing along the error. Otherwise, output holds the successfully opened File.
let output = File::create(filename)?;
```

### A Concurrent Mandelbrot Program

All the pieces are in place, and we can show you the `main` function, where we can put concurrency to work for us. First, a **nonconcurrent** version for simplicity:
```rust
use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() != 5 {
        let program = &args[0];
        eprintln!("Usage: {program} FILE PIXELS LEFT,TOP RIGHT,BOTTOM");
        eprintln!("Example: {program} mandel.png 1000x750 -1.20,0.35 -1,0.20");
        std::process::exit(1);
    }

    let bounds: (usize, usize) = parse_pair(&args[2], 'x')
        .expect("error parsing image dimensions");
    let upper_left = parse_complex(&args[3])
        .expect("error parsing upper left corner point");
    let lower_right = parse_complex(&args[4])
        .expect("error parsing lower right corner point");

    // A macro call vec![v; n] creates a vector n elements long whose elements are
    // initialized to v.
    let mut pixels = vec![0; bounds.0 * bounds.1];

    // The expression &mut pixels borrows a mutable reference to our pixel buffer,
    // allowing render to fill it with computed grayscale values,
    // even while pixels remains the vector’s owner
    render(&mut pixels, bounds, upper_left, lower_right);

    // In this case, we pass a shared (nonmutable) reference to the buffer, since
    // write_image should have no need to modify the buffer’s contents.
    write_image(&args[1], &pixels, bounds)
        .expect("error writing PNG file");
}
```

To this end, we’ll divide the image into sections, one per processor, and let each processor color the pixels assigned to it. For simplicity, we’ll break it into horizontal bands, as shown below. When all processors have finished, we can write out the pixels to disk.
<p align="center"><img src="./assets/maldelbrot-concurrency.png" width="500px" height="auto"></p>

Rust offers a **scoped thread** facility that does exactly what we need here. To use it, we need to take out the single line calling render and replace it with the following:

```rust
// Start by asking the system how many threads we should create.
let threads = std::thread::available_parallelism()
    .expect("error querying CPU count")
    .get();
let rows_per_band = bounds.1.div_ceil(threads);

// Divide the pixel buffer into bands.
let bands = pixels.chunks_mut(rows_per_band * bounds.0);

// std::thread::scope call ensures that all threads have completed before it returns.
std::thread::scope(|spawner| {
    for (i, band) in bands.enumerate() {
        let top = rows_per_band * i;
        let height = band.len() / bounds.0;
        let band_bounds = (bounds.0, height);
        let band_upper_left =
            pixel_to_point(bounds, (0, top), upper_left, lower_right);
        let band_lower_right =
            pixel_to_point(bounds, (bounds.0, top + height),
                           upper_left, lower_right);

        // Finally, we create a thread, running the closure move || { ... }.
        // The move keyword at the front indicates that this closure takes
        // ownership of the variables it uses; in particular, only the
        // closure may use the mutable slice band.
        spawner.spawn(move || {
            render(band, band_bounds, band_upper_left, band_lower_right);
        });
    }
});
```

- Note that, unlike functions declared with `fn`, we don’t need to declare the **types** of a **closure’s arguments**; Rust will infer them, along with its return type
